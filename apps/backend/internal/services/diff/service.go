package diff

import (
	"fmt"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/parser"
	"github.com/storacha/go-ucanto/ucan"
	"github.com/storacha/go-ucanto/validator"
)

type Service struct {
	parser *parser.Service
}

func NewService() *Service {
	return &Service{
		parser: parser.NewService(),
	}
}

func (s *Service) GenerateDiff(parentStr, childStr string) ([]models.CapabilityDiff, error) {
	// 1. Parse both tokens
	parent, err := s.parser.Parse(parentStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parent token: %w", err)
	}

	child, err := s.parser.Parse(childStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse child token: %w", err)
	}

	var diffs []models.CapabilityDiff

	// 2. Iterate through Child Capabilities (What we have now)
	for _, childCap := range child.Capabilities() {
		var matchedParent ucan.Capability[any]
		foundMatch := false

		// 3. Find the Parent capability that authorizes this
		for _, parentCap := range parent.Capabilities() {
			// Leverage go-ucanto's logic to check if Child derives from Parent
			// DefaultDerives checks if 'can' matches and 'with' is a subset
			err := validator.DefaultDerives(childCap, parentCap)
			if err == nil {
				matchedParent = parentCap
				foundMatch = true
				break
			}
		}

		// 4. Determine the Status
		diff := models.CapabilityDiff{
			ChildCap: models.CapToMap(childCap),
		}

		if !foundMatch {
			// ESCALATION: Child has a capability the Parent does not have
			diff.Status = "ADDED"
			diff.Message = "Escalation: Parent does not hold this capability."
		} else {
			diff.ParentCap = models.CapToMap(matchedParent)

			// Check if it was narrowed or identical
			if childCap.With() == matchedParent.With() && childCap.Can() == matchedParent.Can() {
				diff.Status = "UNCHANGED"
				diff.Message = "Capability passed down exactly as is."
			} else {
				diff.Status = "NARROWED"
				diff.Message = fmt.Sprintf("Restricted from '%s' to '%s'", matchedParent.With(), childCap.With())
			}
		}

		diffs = append(diffs, diff)
	}

	return diffs, nil
}