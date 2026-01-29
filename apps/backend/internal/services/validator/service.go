package validator

import (
	"fmt"
	"time"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/parser"
)

type Service struct {
	parser *parser.Service
}

func NewService() *Service {
	return &Service{
		parser: parser.NewService(),
	}
}

// ValidateChain validates a delegation chain
func (s *Service) ValidateChain(tokenBytes []byte) (*models.ValidationResult, error) {
	// 1. Delegate parsing to the Parser Service
	// Since your Parser is fixed, this now works for BOTH CAR files and Raw Tokens!
	chain, err := s.parser.ParseDelegationChain(tokenBytes)
	if err != nil {
		return &models.ValidationResult{
			Valid: false,
			RootCause: &models.ValidationError{
				Type:    "parse_error",
				Message: fmt.Sprintf("Failed to parse UCAN: %v", err),
			},
			Summary: models.ValidationSummary{},
		}, nil
	}

	// 2. Validate the chain links
	var chainLinks []models.ChainLink
	var allIssues []models.ValidationIssue

	for _, del := range chain {
		link := s.validateDelegation(del)
		chainLinks = append(chainLinks, link)
		allIssues = append(allIssues, link.Issues...)
	}

	// 3. Build summary
	summary := s.buildSummary(chainLinks)
	
	// 4. Identify root cause if invalid
	var rootCause *models.ValidationError
	if summary.InvalidLinks > 0 {
		rootCause = s.findRootCause(allIssues, chainLinks[0])
	}

	return &models.ValidationResult{
		Valid:     summary.InvalidLinks == 0,
		Chain:     chainLinks,
		RootCause: rootCause,
		Summary:   summary,
	}, nil
}

// validateDelegation checks a single delegation for issues
func (s *Service) validateDelegation(del *models.DelegationResponse) models.ChainLink {
	var issues []models.ValidationIssue
	now := time.Now()

	// Check 1: Expiration
	if !del.Expiration.IsZero() {
		if del.Expiration.Before(now) {
			timeExpired := now.Sub(del.Expiration)
			issues = append(issues, models.ValidationIssue{
				Type:     "expired",
				Message:  fmt.Sprintf("UCAN expired %v ago", timeExpired.Round(time.Minute)),
				Severity: "error",
			})
		} else if del.Expiration.Before(now.Add(24 * time.Hour)) {
			timeUntilExpiry := del.Expiration.Sub(now)
			issues = append(issues, models.ValidationIssue{
				Type:     "expiring_soon", 
				Message:  fmt.Sprintf("UCAN expires in %v", timeUntilExpiry.Round(time.Minute)),
				Severity: "warning",
			})
		}
	}

	// Check 2: Not Before (nbf)
	if !del.NotBefore.IsZero() && del.NotBefore.After(now) {
		issues = append(issues, models.ValidationIssue{
			Type:     "not_yet_valid",
			Message:  fmt.Sprintf("UCAN not valid until %s", del.NotBefore.Format(time.RFC3339)),
			Severity: "error",
		})
	}

	// Check 3: Capabilities
	if len(del.Capabilities) == 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "no_capabilities",
			Message:  "Delegation has no capabilities",
			Severity: "warning",
		})
	}

	// Check 4: Signature (Optimistic check)
	// If the parser marked it as verified=false but valid=true (raw token), we don't error.
	// But if we had a real signature check failure, we would add an error here.

	// Determine primary capability for display
	var capability models.CapabilityInfo
	if len(del.Capabilities) > 0 {
		capability = del.Capabilities[0]
	}

	valid := s.countErrors(issues) == 0

	return models.ChainLink{
		Level:      del.Level,
		CID:        del.CID,
		Issuer:     del.Issuer,
		Audience:   del.Audience,
		Capability: capability,
		Expiration: del.Expiration,
		NotBefore:  del.NotBefore,
		Valid:      valid,
		Issues:     issues,
	}
}

// Helper: Count severity=error issues
func (s *Service) countErrors(issues []models.ValidationIssue) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == "error" {
			count++
		}
	}
	return count
}

// Helper: Build statistics
func (s *Service) buildSummary(links []models.ChainLink) models.ValidationSummary {
	summary := models.ValidationSummary{
		TotalLinks: len(links),
	}

	for _, link := range links {
		if link.Valid {
			summary.ValidLinks++
		} else {
			summary.InvalidLinks++
		}

		for _, issue := range link.Issues {
			if issue.Severity == "warning" {
				summary.WarningCount++
			}
		}
	}

	return summary
}

// Helper: Find the first error to report as root cause
func (s *Service) findRootCause(issues []models.ValidationIssue, firstLink models.ChainLink) *models.ValidationError {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return &models.ValidationError{
				Type:    issue.Type,
				Message: issue.Message,
				Link: &models.LinkInfo{
					Issuer:   firstLink.Issuer,
					Audience: firstLink.Audience,
				},
			}
		}
	}
	return nil
}