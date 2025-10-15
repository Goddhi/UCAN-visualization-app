package validator

import (
	"fmt"
	"time"

	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/storacha/go-ucanto/ucan"
	"github.com/goddhi/ucan-visualizer/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateChain(tokenBytes []byte) (*models.ValidationResult, error) {
	del, err := delegation.Extract(tokenBytes)
	if err != nil {
		return &models.ValidationResult{
			Valid: false,
			RootCause: &models.ValidationError{
				Type:    "parse_error",
				Message: fmt.Sprintf("Failed to parse UCAN: %v", err),
			},
			Summary: models.ValidationSummary{
				TotalLinks:   0,
				ValidLinks:   0,
				InvalidLinks: 0,
				WarningCount: 0,
			},
		}, nil
	}

	chain := []models.ChainLink{}

	rootLink := s.validateDelegation(del, 0)
	chain = append(chain, rootLink)

	summary := models.ValidationSummary{
		TotalLinks:   len(chain),
		ValidLinks:   0,
		InvalidLinks: 0,
		WarningCount: 0,
	}

	valid := rootLink.Valid
	if valid {
		summary.ValidLinks++
	} else {
		summary.InvalidLinks++
	}

	for _, issue := range rootLink.Issues {
		if issue.Severity == "warning" {
			summary.WarningCount++
		}
	}

	var rootCause *models.ValidationError
	if !valid {
		for _, issue := range rootLink.Issues {
			if issue.Severity == "error" {
				rootCause = &models.ValidationError{
					Type:    issue.Type,
					Message: issue.Message,
					Link: &models.LinkInfo{
						Issuer:   rootLink.Issuer,
						Audience: rootLink.Audience,
					},
				}
				break
			}
		}
	}

	return &models.ValidationResult{
		Valid:     valid,
		Chain:     chain,
		RootCause: rootCause,
		Summary:   summary,
	}, nil
}

// validateDelegation validates a single delegation
func (s *Service) validateDelegation(del delegation.Delegation, level int) models.ChainLink {
	issues := []models.ValidationIssue{}

	now := time.Now()

	exp := del.Expiration()
	var expiration time.Time
	if exp != nil {
		expiration = time.Unix(int64(*exp), 0)

		if ucan.IsExpired(del) {
			timeExpired := now.Sub(expiration)
			issues = append(issues, models.ValidationIssue{
				Type:     "expired",
				Message:  fmt.Sprintf("UCAN expired %v ago at %s", timeExpired.Round(time.Minute), expiration.Format(time.RFC3339)),
				Severity: "error",
				Context: map[string]interface{}{
					"expiration": expiration.Format(time.RFC3339),
					"now":        now.Format(time.RFC3339),
					"expired_by": timeExpired.String(),
				},
			})
		} else if now.Add(24 * time.Hour).After(expiration) {
			timeUntilExpiry := expiration.Sub(now)
			issues = append(issues, models.ValidationIssue{
				Type:     "expiring_soon",
				Message:  fmt.Sprintf("UCAN expires in %v", timeUntilExpiry.Round(time.Minute)),
				Severity: "warning",
				Context: map[string]interface{}{
					"expiration":      expiration.Format(time.RFC3339),
					"time_remaining": timeUntilExpiry.String(),
				},
			})
		}
	} else {
		issues = append(issues, models.ValidationIssue{
			Type:     "no_expiration",
			Message:  "UCAN has no expiration time (valid indefinitely)",
			Severity: "info",
			Context: map[string]interface{}{
				"note": "UCANs without expiration remain valid until revoked",
			},
		})
	}

	nbf := del.NotBefore()
	notBefore := time.Unix(int64(nbf), 0)

	if ucan.IsTooEarly(del) {
		timeUntilValid := notBefore.Sub(now)
		issues = append(issues, models.ValidationIssue{
			Type:     "not_yet_valid",
			Message:  fmt.Sprintf("UCAN not valid until %s (in %v)", notBefore.Format(time.RFC3339), timeUntilValid.Round(time.Minute)),
			Severity: "error",
			Context: map[string]interface{}{
				"not_before":       notBefore.Format(time.RFC3339),
				"now":              now.Format(time.RFC3339),
				"time_until_valid": timeUntilValid.String(),
			},
		})
	}

	if len(del.Proofs()) > 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "has_proofs",
			Message:  fmt.Sprintf("Delegation has %d proof(s) in chain", len(del.Proofs())),
			Severity: "info",
			Context: map[string]interface{}{
				"proof_count": len(del.Proofs()),
				"proof_cids":  s.getProofCIDs(del),
			},
		})
	}

	var capability models.CapabilityInfo
	if len(del.Capabilities()) > 0 {
		cap := del.Capabilities()[0]
		capability = models.CapabilityInfo{
			With: cap.With(),
			Can:  cap.Can(),
			Nb:   make(map[string]interface{}),
		}

		if len(del.Capabilities()) > 1 {
			issues = append(issues, models.ValidationIssue{
				Type:     "multiple_capabilities",
				Message:  fmt.Sprintf("Delegation contains %d capabilities", len(del.Capabilities())),
				Severity: "info",
				Context: map[string]interface{}{
					"capability_count": len(del.Capabilities()),
					"note":             "Only the first capability is displayed in summary",
				},
			})
		}
	} else {
		issues = append(issues, models.ValidationIssue{
			Type:     "no_capabilities",
			Message:  "Delegation has no capabilities",
			Severity: "warning",
			Context: map[string]interface{}{
				"note": "A UCAN should typically contain at least one capability",
			},
		})
	}

	if del.Nonce() != "" {
		issues = append(issues, models.ValidationIssue{
			Type:     "has_nonce",
			Message:  "Delegation includes a nonce",
			Severity: "info",
			Context: map[string]interface{}{
				"nonce": string(del.Nonce()),
			},
		})
	}

	if len(del.Facts()) > 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "has_facts",
			Message:  fmt.Sprintf("Delegation includes %d fact(s)", len(del.Facts())),
			Severity: "info",
			Context: map[string]interface{}{
				"fact_count": len(del.Facts()),
			},
		})
	}

	valid := len(s.filterErrors(issues)) == 0

	return models.ChainLink{
		Level:      level,
		CID:        del.Link().String(),
		Issuer:     del.Issuer().DID().String(),
		Audience:   del.Audience().DID().String(),
		Capability: capability,
		Expiration: expiration,
		NotBefore:  notBefore,
		Valid:      valid,
		Issues:     issues,
	}
}

func (s *Service) getProofCIDs(del delegation.Delegation) []string {
	cids := []string{}
	for _, proof := range del.Proofs() {
		cids = append(cids, proof.String())
	}
	return cids
}

func (s *Service) filterErrors(issues []models.ValidationIssue) []models.ValidationIssue {
	errors := []models.ValidationIssue{}
	for _, issue := range issues {
		if issue.Severity == "error" {
			errors = append(errors, issue)
		}
	}
	return errors
}