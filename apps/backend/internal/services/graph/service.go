package graph

import (
	"fmt"

	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/pkg/utils"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GenerateDelegationGraph(tokenBytes []byte) (*models.GraphResponse, error) {
	// Extract delegation
	del, err := delegation.Extract(tokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract delegation: %w", err)
	}

	// Build graph structure
	nodes := make(map[string]*models.GraphNode)
	edges := []models.GraphEdge{}

	// Build graph from delegation
	s.buildGraph(del, nodes, &edges)

	// Convert nodes map to slice
	nodeSlice := []models.GraphNode{}
	for _, node := range nodes {
		nodeSlice = append(nodeSlice, *node)
	}

	return &models.GraphResponse{
		Nodes: nodeSlice,
		Edges: edges,
	}, nil
}

// buildGraph constructs nodes and edges from delegation
func (s *Service) buildGraph(del delegation.Delegation, nodes map[string]*models.GraphNode, edges *[]models.GraphEdge) {
	// Extract principals
	issuerDID := del.Issuer().DID().String()
	audienceDID := del.Audience().DID().String()

	// Add issuer node (delegator)
	if _, exists := nodes[issuerDID]; !exists {
		nodes[issuerDID] = &models.GraphNode{
			ID:    issuerDID,
			Label: utils.ShortenDID(issuerDID),
			Type:  "root",
			Metadata: map[string]interface{}{
				"fullDid": issuerDID,
				"role":    "issuer",
			},
		}
	}

	// Add audience node (delegatee)
	if _, exists := nodes[audienceDID]; !exists {
		nodes[audienceDID] = &models.GraphNode{
			ID:    audienceDID,
			Label: utils.ShortenDID(audienceDID),
			Type:  "leaf",
			Metadata: map[string]interface{}{
				"fullDid": audienceDID,
				"role":    "audience",
			},
		}
	}

	// Create edges for each capability
	for i, cap := range del.Capabilities() {
		with := cap.With()
		can := cap.Can()

		// Extract caveats if present
		caveats := make(map[string]interface{})
		if cap.Nb() != nil {
			// Try to extract caveats (simplified)
			caveats = s.extractCaveatsSummary(cap.Nb())
		}

		edge := models.GraphEdge{
			Source: issuerDID,
			Target: audienceDID,
			Capability: models.CapabilityInfo{
				With: with,
				Can:  can,
				Nb:   caveats,
			},
			Valid: true,
			Label: s.createEdgeLabel(can, with, i),
		}
		*edges = append(*edges, edge)
	}

	// Add proof metadata
	if len(del.Proofs()) > 0 {
		if node, exists := nodes[audienceDID]; exists {
			node.Metadata["hasProofs"] = true
			node.Metadata["proofCount"] = len(del.Proofs())

			// Store proof CIDs
			proofCIDs := []string{}
			for _, proofLink := range del.Proofs() {
				proofCIDs = append(proofCIDs, proofLink.String())
			}
			node.Metadata["proofCIDs"] = proofCIDs
		}
	}

	// Add expiration metadata
	exp := del.Expiration()
	if exp != nil {
		if node, exists := nodes[audienceDID]; exists {
			node.Metadata["expiration"] = *exp
		}
	}

	// Add not-before metadata
	nbf := del.NotBefore()
	if nbf != 0 {
		if node, exists := nodes[audienceDID]; exists {
			node.Metadata["notBefore"] = nbf
		}
	}

	// Add CID to audience node
	if node, exists := nodes[audienceDID]; exists {
		node.Metadata["delegationCID"] = del.Link().String()
	}
}

// createEdgeLabel creates a human-readable label for the edge
func (s *Service) createEdgeLabel(can, with string, index int) string {
	// If there are multiple capabilities, add an index
	label := fmt.Sprintf("%s", can)

	// Add resource info if it's not too long
	if len(with) < 40 {
		label = fmt.Sprintf("%s on %s", can, with)
	}

	// For multiple capabilities, add numbering
	if index > 0 {
		label = fmt.Sprintf("[%d] %s", index+1, label)
	}

	return label
}

// extractCaveatsSummary extracts a summary of caveats for display
func (s *Service) extractCaveatsSummary(nb any) map[string]interface{} {
	result := make(map[string]interface{})

	if nb == nil {
		return result
	}

	// Simple summary - just indicate there are caveats
	// Full extraction would require ipld.Node handling (like in parser)
	result["_hasCaveats"] = true

	return result
}