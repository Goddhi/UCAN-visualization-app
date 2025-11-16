package validator

import (
	"fmt"
	"regexp"
	"strings"
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

func (s *Service) ValidateChain(tokenBytes []byte) (*models.ValidationResult, error) {
	// Use parser service to get the chain
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

	// Validate each delegation in the chain
	var chainLinks []models.ChainLink
	var allIssues []models.ValidationIssue

	for _, del := range chain {
		link := s.validateDelegation(del)
		chainLinks = append(chainLinks, link)
		allIssues = append(allIssues, link.Issues...)
	}

	// Build summary
	summary := s.buildSummary(chainLinks)
	
	// Find root cause if invalid
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

func (s *Service) validateDelegation(del *models.DelegationResponse) models.ChainLink {
	var issues []models.ValidationIssue
	now := time.Now()

	// Check expiration
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

	// Check not-before
	if del.NotBefore.After(now) {
		issues = append(issues, models.ValidationIssue{
			Type:     "not_yet_valid",
			Message:  fmt.Sprintf("UCAN not valid until %s", del.NotBefore.Format(time.RFC3339)),
			Severity: "error",
		})
	}

	// Check capabilities
	if len(del.Capabilities) == 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "no_capabilities",
			Message:  "Delegation has no capabilities",
			Severity: "warning",
		})
	}

	// Check proofs (basic info)
	if len(del.Proofs) > 0 {
		issues = append(issues, models.ValidationIssue{
			Type:     "has_proofs",
			Message:  fmt.Sprintf("Delegation has %d proof(s)", len(del.Proofs)),
			Severity: "info",
		})
	}

	// Determine primary capability
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

func (s *Service) validateCapabilityAttenuation(parent, child models.CapabilityInfo) []models.ValidationIssue {
	var issues []models.ValidationIssue

	// Check resource (with field)
	if !s.resourceMatches(parent.With, child.With) {
		issues = append(issues, models.ValidationIssue{
			Type:     "resource_mismatch",
			Message:  fmt.Sprintf("Child resource '%s' not covered by parent '%s'", child.With, parent.With),
			Severity: "error",
		})
	}

	// Check ability (can field)
	if !s.abilityMatches(parent.Can, child.Can) {
		issues = append(issues, models.ValidationIssue{
			Type:     "capability_escalation",
			Message:  fmt.Sprintf("Child ability '%s' exceeds parent '%s'", child.Can, parent.Can),
			Severity: "error",
		})
	}

	return issues
}

func (s *Service) resourceMatches(parent, child string) bool {
	if parent == child {
		return true
	}

	// Wildcard matching: "storage:*" matches "storage:alice/*"
	pattern := strings.ReplaceAll(regexp.QuoteMeta(parent), `\*`, ".*")
	pattern = "^" + pattern + "$"
	matched, _ := regexp.MatchString(pattern, child)
	return matched
}

func (s *Service) abilityMatches(parent, child string) bool {
	if parent == child || parent == "*" {
		return true
	}

	// Wildcard matching: "store/*" matches "store/add"
	if strings.HasSuffix(parent, "/*") {
		prefix := strings.TrimSuffix(parent, "/*")
		return strings.HasPrefix(child, prefix+"/")
	}

	return false
}

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

func (s *Service) countErrors(issues []models.ValidationIssue) int {
	count := 0
	for _, issue := range issues {
		if issue.Severity == "error" {
			count++
		}
	}
	return count
}