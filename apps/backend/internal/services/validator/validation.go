package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/goddhi/ucan-visualizer/internal/models"
)

// validateCapabilityDelegation checks if child capability is properly attenuated
func (s *Service) validateCapabilityDelegation(parent, child models.CapabilityInfo) []models.ValidationIssue {
	issues := []models.ValidationIssue{}

	// 1. Check resource (with field)
	if !s.resourceMatches(parent.With, child.With) {
		issues = append(issues, models.ValidationIssue{
			Type:     "resource_mismatch",
			Message:  fmt.Sprintf("Child resource '%s' not covered by parent '%s'", child.With, parent.With),
			Severity: "error",
			Context: map[string]interface{}{
				"parent": parent.With,
				"child":  child.With,
			},
		})
	}

	// 2. Check ability (can field)
	if !s.abilityMatches(parent.Can, child.Can) {
		issues = append(issues, models.ValidationIssue{
			Type:     "capability_escalation",
			Message:  fmt.Sprintf("Child ability '%s' exceeds parent '%s'", child.Can, parent.Can),
			Severity: "error",
			Context: map[string]interface{}{
				"parent": parent.Can,
				"child":  child.Can,
			},
		})
	}

	// 3. Check caveats (nb field)
	caveatIssues := s.validateCaveats(parent.Nb, child.Nb)
	issues = append(issues, caveatIssues...)

	return issues
}

// resourceMatches checks if child resource is covered by parent
func (s *Service) resourceMatches(parent, child string) bool {
	// Exact match
	if parent == child {
		return true
	}

	// Wildcard matching
	// Examples:
	//   parent: "storage:*" matches child: "storage:alice/*"
	//   parent: "storage:alice/*" matches child: "storage:alice/photos"
	//   parent: "storage:alice/*" does NOT match child: "storage:bob/*"

	pattern := strings.ReplaceAll(regexp.QuoteMeta(parent), `\*`, ".*")
	pattern = "^" + pattern + "$"
	matched, _ := regexp.MatchString(pattern, child)

	return matched
}

// abilityMatches checks if child ability is covered by parent
func (s *Service) abilityMatches(parent, child string) bool {
	// Exact match
	if parent == child {
		return true
	}

	// Wildcard matching
	// Examples:
	//   parent: "store/*" matches child: "store/add"
	//   parent: "store/add" matches child: "store/add"
	//   parent: "store/add" does NOT match child: "store/remove"

	if strings.HasSuffix(parent, "/*") {
		prefix := strings.TrimSuffix(parent, "/*")
		return strings.HasPrefix(child, prefix+"/")
	}

	// Universal wildcard
	if parent == "*" {
		return true
	}

	return false
}

// validateCaveats ensures child caveats are equally or more restrictive
func (s *Service) validateCaveats(parentNb, childNb map[string]interface{}) []models.ValidationIssue {
	issues := []models.ValidationIssue{}

	// Check numeric restrictions (e.g., size limits)
	if parentSize, ok := parentNb["size"].(float64); ok {
		if childSize, ok := childNb["size"].(float64); ok {
			if childSize > parentSize {
				issues = append(issues, models.ValidationIssue{
					Type:     "caveat_escalation",
					Message:  fmt.Sprintf("Child size limit (%.0f) exceeds parent (%.0f)", childSize, parentSize),
					Severity: "error",
				})
			}
		}
		// If child doesn't have size caveat but parent does, that's okay (more restrictive)
	}

	// Check array restrictions (e.g., allowed content types)
	if parentTypes, ok := parentNb["contentType"].([]interface{}); ok {
		if childTypes, ok := childNb["contentType"].([]interface{}); ok {
			// Child types must be subset of parent types
			if !isSubset(childTypes, parentTypes) {
				issues = append(issues, models.ValidationIssue{
					Type:     "caveat_escalation",
					Message:  "Child content types exceed parent restrictions",
					Severity: "error",
				})
			}
		}
	}

	return issues
}

// isSubset checks if all elements in child are in parent
func isSubset(child, parent []interface{}) bool {
	parentSet := make(map[interface{}]bool)
	for _, item := range parent {
		parentSet[item] = true
	}

	for _, item := range child {
		if !parentSet[item] {
			return false
		}
	}

	return true
}