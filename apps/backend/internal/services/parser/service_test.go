package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/goddhi/ucan-visualizer/test/fixtures"
)

func TestParseDelegation(t *testing.T) {
	svc := NewService()

	t.Run("Successfully parses valid CAR token", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateValidUCAN()
		require.NoError(t, err)

		result, err := svc.ParseDelegation(tokenBytes)

		require.NoError(t, err)
		assert.NotEmpty(t, result.CID)
		assert.NotEmpty(t, result.Issuer)
		assert.NotEmpty(t, result.Audience)
		
		assert.Greater(t, len(result.Capabilities), 0)
		assert.Equal(t, "store/add", result.Capabilities[0].Can)
	})

	t.Run("Fails on garbage input", func(t *testing.T) {
		garbageBytes := []byte("this is not a real UCAN")

		result, err := svc.ParseDelegation(garbageBytes)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode")
	})
}

func TestParseDelegationChain(t *testing.T) {
	svc := NewService()

	t.Run("Parses multi-hop chain", func(t *testing.T) {
		chainBytes, err := fixtures.GenerateComplexChain()
		require.NoError(t, err)

		chain, err := svc.ParseDelegationChain(chainBytes)

		require.NoError(t, err)
		assert.NotEmpty(t, chain)
	})
}