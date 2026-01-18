package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/goddhi/ucan-visualizer/test/fixtures"
)

func TestGenerateDelegationGraph(t *testing.T) {
	svc := NewService()

	t.Run("Generates nodes and edges for chain", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateComplexChain()
		require.NoError(t, err)

		graph, err := svc.GenerateDelegationGraph(tokenBytes)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(graph.Nodes), 2)
		
		hasRoot := false
		hasLeaf := false
		for _, node := range graph.Nodes {
			if node.Type == "root" { hasRoot = true }
			if node.Type == "leaf" { hasLeaf = true }
		}
		assert.True(t, hasRoot, "Graph must have a root node")
		assert.True(t, hasLeaf, "Graph must have a leaf node")

		assert.NotEmpty(t, graph.Edges)
		assert.Equal(t, "delegation", graph.Edges[0].Type)
	})
}

func TestGenerateInvocationGraph(t *testing.T) {
	svc := NewService()

	t.Run("Identifies delegation as NOT invocation", func(t *testing.T) {
		tokenBytes, err := fixtures.GenerateValidUCAN()
		require.NoError(t, err)

		graph, err := svc.GenerateInvocationGraph(tokenBytes)
		require.NoError(t, err)

		assert.False(t, graph.IsInvocation)
		assert.NotEmpty(t, graph.Nodes)
	})
}