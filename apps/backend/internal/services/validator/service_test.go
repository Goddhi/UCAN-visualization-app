package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/goddhi/ucan-visualizer/test/fixtures"
)

func TestValidateChain(t *testing.T) {
	svc := NewService()

	t.Run("Validates chain", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateComplexChain()
		require.NoError(t, err)

		result, err := svc.ValidateChain(tokenBytes)
		require.NoError(t, err)

		assert.True(t, result.Valid)
		assert.Equal(t, 0, result.Summary.InvalidLinks)
		assert.Nil(t, result.RootCause)
		
		assert.GreaterOrEqual(t, len(result.Chain), 1)
	})

	t.Run("Detects expired tokens", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateExpiredUCAN()
		require.NoError(t, err)

		result, err := svc.ValidateChain(tokenBytes)
		require.NoError(t, err)

		assert.False(t, result.Valid)
		assert.NotNil(t, result.RootCause)
		assert.Equal(t, "expired", result.RootCause.Type)
		assert.Contains(t, result.RootCause.Message, "UCAN expired")
	})
}

func TestResourceMatching(t *testing.T) {
	svc := NewService()

	tests := []struct {
		parent   string
		child    string
		expected bool
	}{
		{"storage:*", "storage:alice/photos", true},
		{"storage:alice/*", "storage:alice/photos", true},
		{"storage:alice/*", "storage:bob/photos", false},
		{"*", "anything", true},
	}

	for _, tt := range tests {
		t.Run(tt.parent+" -> "+tt.child, func(t *testing.T) {
			assert.Equal(t, tt.expected, svc.resourceMatches(tt.parent, tt.child))
		})
	}
}