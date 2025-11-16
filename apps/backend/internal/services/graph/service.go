package graph

import (
	"fmt"
	"sort"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/parser"
	"github.com/goddhi/ucan-visualizer/pkg/utils"
)

type Service struct {
	parser *parser.Service
}

func NewService() *Service {
	return &Service{
		parser: parser.NewService(),
	}
}

// GenerateDelegationGraph creates comprehensive delegation chain visualization
func (s *Service) GenerateDelegationGraph(tokenBytes []byte) (*models.GraphResponse, error) {
	chain, err := s.parser.ParseDelegationChain(tokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse delegation chain: %w", err)
	}

	nodes, edges := s.buildDelegationGraph(chain)
	chainInfo := s.buildChainInfo(chain)

	return &models.GraphResponse{
		Nodes: nodes,
		Edges: edges,
		Chain: chainInfo,
	}, nil
}

// GenerateInvocationGraph creates enhanced invocation visualization
func (s *Service) GenerateInvocationGraph(tokenBytes []byte) (*models.InvocationGraphResponse, error) {
	invocation, err := s.parser.ParseInvocation(tokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse invocation: %w", err)
	}

	chain, err := s.parser.ParseDelegationChain(tokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse delegation chain: %w", err)
	}

	nodes, edges := s.buildInvocationGraph(chain, invocation)
	chainInfo := s.buildChainInfo(chain)

	return &models.InvocationGraphResponse{
		Nodes:        nodes,
		Edges:        edges,
		Chain:        chainInfo,
		Invocation:   invocation,
		IsInvocation: invocation.IsInvocation,
	}, nil
}

// buildDelegationGraph creates nodes and edges for delegation visualization
func (s *Service) buildDelegationGraph(chain []*models.DelegationResponse) ([]models.GraphNode, []models.GraphEdge) {
	nodes := make(map[string]*models.GraphNode)
	var edges []models.GraphEdge

	// Find max level for proper node typing
	maxLevel := 0
	for _, del := range chain {
		if del.Level > maxLevel {
			maxLevel = del.Level
		}
	}

	// Create nodes and edges for each delegation
	for _, del := range chain {
		// Create issuer node
		if _, exists := nodes[del.Issuer]; !exists {
			nodes[del.Issuer] = &models.GraphNode{
				ID:    del.Issuer,
				Label: utils.ShortenDID(del.Issuer),
				Type:  s.getNodeType(del.Level, maxLevel, "issuer"),
				Level: del.Level,
				Metadata: map[string]interface{}{
					"fullDID":       del.Issuer,
					"role":          "delegator",
					"capabilities":  len(del.Capabilities),
					"proofs":        len(del.Proofs),
				},
			}
		}

		// Create audience node
		if _, exists := nodes[del.Audience]; !exists {
			nodes[del.Audience] = &models.GraphNode{
				ID:    del.Audience,
				Label: utils.ShortenDID(del.Audience),
				Type:  s.getNodeType(del.Level, maxLevel, "audience"),
				Level: del.Level,
				Metadata: map[string]interface{}{
					"fullDID": del.Audience,
					"role":    "delegatee",
				},
			}
		}

		// Create edges for each capability
		for i, cap := range del.Capabilities {
			edge := models.GraphEdge{
				Source:     del.Issuer,
				Target:     del.Audience,
				Capability: cap,
				Label:      s.createCapabilityLabel(cap, i),
				Valid:      true,
				Level:      del.Level,
				Type:       "delegation",
				Metadata: map[string]interface{}{
					"capability": cap.Can,
					"resource":   cap.With,
					"category":   cap.Category,
					"cid":        del.CID,
				},
			}
			edges = append(edges, edge)
		}

		// Add proof edges if they exist
		for _, proof := range del.Proofs {
			// Create proof connection edges
			proofEdge := models.GraphEdge{
				Source: fmt.Sprintf("proof-%s", proof.CID),
				Target: del.Issuer,
				Label:  fmt.Sprintf("Proof %d", proof.Index),
				Valid:  true,
				Level:  del.Level + 1,
				Type:   "proof",
				Metadata: map[string]interface{}{
					"proofCID":  proof.CID,
					"proofType": proof.Type,
				},
			}
			edges = append(edges, proofEdge)
		}
	}

	// Convert nodes map to slice
	var nodeSlice []models.GraphNode
	for _, node := range nodes {
		nodeSlice = append(nodeSlice, *node)
	}

	return nodeSlice, edges
}

// buildInvocationGraph creates enhanced visualization for invocations
func (s *Service) buildInvocationGraph(chain []*models.DelegationResponse, invocation *models.InvocationResponse) ([]models.GraphNode, []models.GraphEdge) {
	nodes, edges := s.buildDelegationGraph(chain)

	// Enhance visualization for invocations
	if invocation.IsInvocation && invocation.Task != nil {
		// Mark invoker node
		for i, node := range nodes {
			if node.ID == invocation.Task.Issuer {
				nodes[i].Type = "invoker"
				nodes[i].Metadata["isInvoker"] = true
				nodes[i].Metadata["taskAction"] = invocation.Task.Action
				nodes[i].Metadata["taskType"] = invocation.Task.TaskType
				nodes[i].Metadata["permissions"] = invocation.Task.Permissions
			}
			if node.ID == invocation.Task.Target {
				nodes[i].Metadata["isTarget"] = true
				nodes[i].Metadata["targetResource"] = invocation.Task.Resource
			}
		}

		// Mark invocation edges
		for i, edge := range edges {
			if edge.Source == invocation.Task.Issuer && edge.Target == invocation.Task.Target {
				edges[i].Type = "invocation"
				edges[i].Label = fmt.Sprintf("INVOKE: %s", edge.Label)
				edges[i].Metadata["isInvocation"] = true
				edges[i].Metadata["taskAction"] = invocation.Task.Action
				edges[i].Metadata["constraints"] = invocation.Task.Constraints

				// Add invocation analysis metadata
				if invocation.InvocationAnalysis != nil {
					edges[i].Metadata["hasInvokeCapability"] = invocation.InvocationAnalysis.HasInvokeCapability
					edges[i].Metadata["invokePatterns"] = invocation.InvocationAnalysis.InvokePatterns
				}
			}
		}

		// Add capability analysis nodes if significant
		if invocation.CapabilityAnalysis != nil && len(invocation.CapabilityAnalysis.Categories) > 1 {
			// Create capability summary node
			capNode := models.GraphNode{
				ID:    "capability-summary",
				Label: fmt.Sprintf("Capabilities (%d)", invocation.CapabilityAnalysis.TotalCount),
				Type:  "summary",
				Level: -1, // Special level for summary nodes
				Metadata: map[string]interface{}{
					"totalCapabilities": invocation.CapabilityAnalysis.TotalCount,
					"invokeCount":       invocation.CapabilityAnalysis.InvokeCount,
					"delegateCount":     invocation.CapabilityAnalysis.DelegateCount,
					"categories":        s.getCategoryNames(invocation.CapabilityAnalysis.Categories),
				},
			}
			nodes = append(nodes, capNode)
		}
	}

	return nodes, edges
}

// buildChainInfo creates comprehensive chain analysis
func (s *Service) buildChainInfo(chain []*models.DelegationResponse) models.ChainInfo {
	if len(chain) == 0 {
		return models.ChainInfo{}
	}

	principals := make(map[string]*models.PrincipalInfo)
	var timeline []models.TimelineEvent
	maxLevel := 0
	var proofChain *models.ProofChain

	// Find max level and build proof chain
	for _, del := range chain {
		if del.Level > maxLevel {
			maxLevel = del.Level
		}
	}

	// Build proof chain structure
	if maxLevel > 0 || s.hasProofs(chain) {
		proofChain = s.buildProofChain(chain)
	}

	// Process each delegation
	for _, del := range chain {
		// Track principals
		for _, did := range []string{del.Issuer, del.Audience} {
			if _, exists := principals[did]; !exists {
				role := s.getPrincipalRole(did, del.Level, maxLevel)
				principals[did] = &models.PrincipalInfo{
					DID:   did,
					Role:  role,
					Level: del.Level,
					CIDs:  []string{del.CID},
				}
			} else {
				principals[did].CIDs = append(principals[did].CIDs, del.CID)
			}
		}

		// Add timeline events
		if !del.NotBefore.IsZero() {
			timeline = append(timeline, models.TimelineEvent{
				Type:      "issued",
				Time:      del.NotBefore,
				CID:       del.CID,
				Level:     del.Level,
				Principal: del.Issuer,
			})
		}

		if !del.Expiration.IsZero() {
			timeline = append(timeline, models.TimelineEvent{
				Type:      "expires",
				Time:      del.Expiration,
				CID:       del.CID,
				Level:     del.Level,
				Principal: del.Audience,
			})
		}
	}

	// Sort timeline
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Time.Before(timeline[j].Time)
	})

	// Convert principals to slice
	var principalSlice []models.PrincipalInfo
	for _, principal := range principals {
		principalSlice = append(principalSlice, *principal)
	}

	// Get leaf CIDs
	var leafCIDs []string
	for _, del := range chain {
		if del.Level == maxLevel {
			leafCIDs = append(leafCIDs, del.CID)
		}
	}

	return models.ChainInfo{
		TotalLevels: maxLevel + 1,
		IsComplete:  true,
		RootCID:     chain[0].CID,
		LeafCIDs:    leafCIDs,
		Principals:  principalSlice,
		Timeline:    timeline,
		ProofChain:  proofChain,
	}
}

// Helper functions
func (s *Service) getNodeType(level, maxLevel int, role string) string {
	switch {
	case level == 0 && role == "issuer":
		return "root"
	case level == maxLevel:
		return "leaf"
	case level > 0:
		return "intermediate"
	default:
		return "node"
	}
}

func (s *Service) createCapabilityLabel(cap models.CapabilityInfo, index int) string {
	label := cap.Can
	
	// Add resource if not too long
	if len(cap.With) < 30 {
		label = fmt.Sprintf("%s on %s", cap.Can, cap.With)
	}
	
	// Add capability category
	if cap.Category != "" && cap.Category != "general" {
		label = fmt.Sprintf("[%s] %s", cap.Category, label)
	}
	
	// Add index for multiple capabilities
	if index > 0 {
		label = fmt.Sprintf("%d. %s", index+1, label)
	}
	
	return label
}

func (s *Service) getPrincipalRole(did string, level, maxLevel int) string {
	switch {
	case level == 0:
		return "root"
	case level == maxLevel:
		return "leaf"
	default:
		return "intermediate"
	}
}

func (s *Service) hasProofs(chain []*models.DelegationResponse) bool {
	for _, del := range chain {
		if len(del.Proofs) > 0 {
			return true
		}
	}
	return false
}

func (s *Service) buildProofChain(chain []*models.DelegationResponse) *models.ProofChain {
	levels := make(map[int]*models.ProofLevel)
	
	for _, del := range chain {
		if _, exists := levels[del.Level]; !exists {
			levels[del.Level] = &models.ProofLevel{
				Level:       del.Level,
				Delegations: []*models.DelegationResponse{},
				ProofLinks:  []string{},
			}
		}
		
		levels[del.Level].Delegations = append(levels[del.Level].Delegations, del)
		
		for _, proof := range del.Proofs {
			levels[del.Level].ProofLinks = append(levels[del.Level].ProofLinks, proof.CID)
		}
	}
	
	// Convert to slice
	var levelSlice []models.ProofLevel
	for _, level := range levels {
		levelSlice = append(levelSlice, *level)
	}
	
	// Sort by level
	sort.Slice(levelSlice, func(i, j int) bool {
		return levelSlice[i].Level < levelSlice[j].Level
	})
	
	return &models.ProofChain{
		Root:       chain[0].CID,
		Levels:     levelSlice,
		IsComplete: true,
		TotalDepth: len(levelSlice),
	}
}

func (s *Service) getCategoryNames(categories map[string][]models.CapabilityInfo) []string {
	var names []string
	for category := range categories {
		names = append(names, category)
	}
	return names
}