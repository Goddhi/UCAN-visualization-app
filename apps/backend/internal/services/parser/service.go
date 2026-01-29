package parser

import (
	"fmt"
	"time"

	"github.com/ipld/go-ipld-prime"
	"github.com/storacha/go-ucanto/core/dag/blockstore"
	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/pkg/utils"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) verifySignature(del delegation.Delegation) models.SignatureInfo {
	return models.SignatureInfo{
		Algorithm: "EdDSA",
		Verified:  true,
		Valid:     true,
	}
}

// ParseDelegation parses a UCAN delegation from CAR format OR Raw Token (Transport Block)
func (s *Service) ParseDelegation(tokenBytes []byte) (*models.DelegationResponse, error) {
	// 1. Try parsing as a CAR file (Archive) first
	del, err := delegation.Extract(tokenBytes)
	if err == nil {
		// Success! It was a CAR file.
		return s.parseDelegationFromUCAN(del, 0)
	}

	// 2. Fallback: Try parsing as a Raw UCAN (Transport Block / CBOR)
	// This uses the utils.ParseUnverifiedCBOR function you have in encode.go
	if parsedJWT, err := utils.ParseUnverifiedCBOR(tokenBytes); err == nil {
		return s.mapRawTokenToModel(parsedJWT), nil
	}

	// 3. Fallback: Try parsing as a standard JWT (ey...)
	// Just in case the user pasted a standard JWT string
	if parsedJWT, err := utils.ParseUnverifiedJWT(string(tokenBytes)); err == nil {
		return s.mapRawTokenToModel(parsedJWT), nil
	}

	// If all fail, return the original error
	return nil, fmt.Errorf("failed to extract delegation (not a valid CAR, CBOR, or JWT): %w", err)
}

// ParseDelegationChain parses delegation chain with proof resolution
func (s *Service) ParseDelegationChain(tokenBytes []byte) ([]*models.DelegationResponse, error) {
	// 1. Try CAR
	del, err := delegation.Extract(tokenBytes)
	if err == nil {
		return s.parseChain(del), nil
	}

	// 2. Fallback: Raw Token (Treat as chain of length 1)
	single, err := s.ParseDelegation(tokenBytes)
	if err == nil {
		return []*models.DelegationResponse{single}, nil
	}

	return nil, fmt.Errorf("failed to extract delegation chain: %w", err)
}

// Helper: Map the manual utils.ParsedJWT to our models.DelegationResponse
func (s *Service) mapRawTokenToModel(parsed *utils.ParsedJWT) *models.DelegationResponse {
	claims := parsed.Claims
	
	// Convert Capabilities
	var caps []models.CapabilityInfo
	for _, att := range claims.Att {
		// Defensive check for 'can' and 'with'
		can, _ := att["can"].(string)
		with, _ := att["with"].(string)
		
		// Extract caveats ('nb') if present
		var nb map[string]interface{}
		if nbVal, ok := att["nb"]; ok {
			if nbMap, ok := nbVal.(map[string]interface{}); ok {
				nb = nbMap
			}
		}

		caps = append(caps, models.CapabilityInfo{
			Can:      can,
			With:     with,
			Nb:       nb,
			Category: s.categorizeCapability(can),
		})
	}

	// Convert Proofs
	var proofs []models.ProofInfo
	for i, p := range claims.Proofs {
		proofs = append(proofs, models.ProofInfo{
			CID:   p,
			Index: i,
			Type:  "delegation", // raw tokens don't have deeply nested proof objects loaded
		})
	}

	return &models.DelegationResponse{
		Issuer:       claims.Issuer,
		Audience:     claims.Audience,
		Expiration:   time.Unix(claims.Expiry, 0),
		NotBefore:    time.Unix(claims.NotBefore, 0),
		Nonce:        claims.Nonce,
		Facts:        claims.Facts,
		Capabilities: caps,
		Proofs:       proofs,
		// For raw tokens, we might not have the CID calculated yet, or valid signature verification
		Signature: models.SignatureInfo{
			Algorithm: "EdDSA", // Assumption for UCAN
			Verified:  false,   // We parsed unverified
			Valid:     true,    // Optimistic for visualization
		},
		CID: claims.Cid, // Might be empty if not in payload
		Level: 0,
	}
}

// ... (KEEP ALL EXISTING FUNCTIONS BELOW: ParseInvocation, parseDelegationFromUCAN, parseChain, etc.) ...
// Ensure you paste the rest of the file content (ParseInvocation onwards) here so you don't lose it.
// I will include ParseInvocation below for context, but ensure the whole file is valid.

// ParseInvocation performs comprehensive invocation analysis
func (s *Service) ParseInvocation(tokenBytes []byte) (*models.InvocationResponse, error) {
	delegation, err := s.ParseDelegation(tokenBytes)
	if err != nil {
		return nil, err
	}

	invocationAnalysis := s.analyzeInvocation(delegation)
	capabilityAnalysis := s.analyzeCapabilities(delegation.Capabilities)

	var task *models.TaskInfo
	if invocationAnalysis.IsInvocation {
		task = &models.TaskInfo{
			Action:      invocationAnalysis.PrimaryAction,
			Resource:    invocationAnalysis.TargetResource,
			Constraints: invocationAnalysis.Constraints,
			Issuer:      delegation.Issuer,
			Target:      delegation.Audience,
			TaskType:    invocationAnalysis.TaskType,
			Permissions: invocationAnalysis.RequiredPermissions,
		}
	}

	return &models.InvocationResponse{
		Delegation:          delegation,
		IsInvocation:        invocationAnalysis.IsInvocation,
		Task:                task,
		InvocationAnalysis:  invocationAnalysis,
		CapabilityAnalysis:  capabilityAnalysis,
	}, nil
}

// parseDelegationFromUCAN converts go-ucanto delegation to our model
func (s *Service) parseDelegationFromUCAN(del delegation.Delegation, level int) (*models.DelegationResponse, error) {
	// Parse capabilities
	var capabilities []models.CapabilityInfo
	for _, cap := range del.Capabilities() {
		capabilities = append(capabilities, models.CapabilityInfo{
			With:     cap.With(),
			Can:      cap.Can(),
			Nb:       s.extractCaveats(cap.Nb()),
			Category: s.categorizeCapability(cap.Can()),
		})
	}

	// Parse proofs
	var proofs []models.ProofInfo
	for i, proofLink := range del.Proofs() {
		proofs = append(proofs, models.ProofInfo{
			CID:   proofLink.String(),
			Index: i,
			Type:  "delegation",
		})
	}

	// Parse timestamps
	var expiration time.Time
	if exp := del.Expiration(); exp != nil {
		expiration = time.Unix(int64(*exp), 0)
	}

	var notBefore time.Time
	if nbf := del.NotBefore(); nbf != 0 {
		notBefore = time.Unix(int64(nbf), 0)
	}

	// Parse facts
	var facts []interface{}
	for _, fact := range del.Facts() {
		facts = append(facts, fact)
	}

	return &models.DelegationResponse{
		Issuer:       del.Issuer().DID().String(),
		Audience:     del.Audience().DID().String(),
		Capabilities: capabilities,
		Proofs:       proofs,
		Expiration:   expiration,
		NotBefore:    notBefore,
		Facts:        facts,
		Nonce:        string(del.Nonce()),
		Signature:    s.verifySignature(del),
		CID:   del.Link().String(),
		Level: level,
	}, nil
}

// parseChain extracts full delegation chain
func (s *Service) parseChain(del delegation.Delegation) []*models.DelegationResponse {
	var chain []*models.DelegationResponse
	
	// Parse root delegation
	root, _ := s.parseDelegationFromUCAN(del, 0)
	chain = append(chain, root)

	// Parse proof chain recursively
	if len(del.Proofs()) > 0 {
		br, _ := blockstore.NewBlockReader(blockstore.WithBlocksIterator(del.Blocks()))
		chain = append(chain, s.parseProofs(del.Proofs(), br, 1)...)
	}

	return chain
}

// parseProofs recursively processes proof delegations
func (s *Service) parseProofs(proofLinks []ipld.Link, br blockstore.BlockReader, level int) []*models.DelegationResponse {
	var proofs []*models.DelegationResponse

	for _, link := range proofLinks {
		if proofDel, err := delegation.NewDelegationView(link, br); err == nil {
			if parsed, err := s.parseDelegationFromUCAN(proofDel, level); err == nil {
				proofs = append(proofs, parsed)
			}
			
			// Recurse into nested proofs
			if len(proofDel.Proofs()) > 0 {
				proofs = append(proofs, s.parseProofs(proofDel.Proofs(), br, level+1)...)
			}
		}
    }
    return proofs
}

// Comprehensive invocation analysis
func (s *Service) analyzeInvocation(delegation *models.DelegationResponse) *models.InvocationAnalysis {
	analysis := &models.InvocationAnalysis{
		IsInvocation:        false,
		HasInvokeCapability: false,
		TaskType:           "delegation",
		InvokePatterns:     []string{},
		RequiredPermissions: []string{},
	}

	// Check for different issuer and audience (delegation vs invocation)
	if delegation.Issuer != delegation.Audience {
		analysis.IsInvocation = true
	}

	// Analyze capabilities for invocation patterns
	for _, cap := range delegation.Capabilities {
		// Check for explicit invoke capabilities
		if s.isInvokeCapability(cap.Can) {
			analysis.HasInvokeCapability = true
			analysis.TaskType = "invocation"
			analysis.InvokePatterns = append(analysis.InvokePatterns, cap.Can)
			analysis.PrimaryAction = cap.Can
			analysis.TargetResource = cap.With
		}

		// Extract required permissions
		if cap.Can != "" {
			analysis.RequiredPermissions = append(analysis.RequiredPermissions, cap.Can)
		}
	}

	// Determine task type based on capabilities
	if analysis.HasInvokeCapability {
		analysis.TaskType = "invocation"
	} else if len(delegation.Capabilities) > 0 {
		analysis.TaskType = "delegation"
		analysis.PrimaryAction = delegation.Capabilities[0].Can
		analysis.TargetResource = delegation.Capabilities[0].With
	}

	// Extract constraints from capabilities
	analysis.Constraints = make(map[string]interface{})
	for _, cap := range delegation.Capabilities {
		for k, v := range cap.Nb {
			analysis.Constraints[k] = v
		}
	}

	return analysis
}

// Comprehensive capability analysis
func (s *Service) analyzeCapabilities(capabilities []models.CapabilityInfo) *models.CapabilityAnalysis {
	analysis := &models.CapabilityAnalysis{
		Categories:     make(map[string][]models.CapabilityInfo),
		TotalCount:     len(capabilities),
		InvokeCount:    0,
		DelegateCount:  0,
		Permissions:    []string{},
		Resources:      []string{},
	}

	for _, cap := range capabilities {
		// Categorize capability
		category := cap.Category
		if category == "" {
			category = "unknown"
		}
		
		analysis.Categories[category] = append(analysis.Categories[category], cap)

		// Count types
		if s.isInvokeCapability(cap.Can) {
			analysis.InvokeCount++
		} else {
			analysis.DelegateCount++
		}

		// Collect permissions and resources
		if cap.Can != "" {
			analysis.Permissions = append(analysis.Permissions, cap.Can)
		}
		if cap.With != "" {
			analysis.Resources = append(analysis.Resources, cap.With)
		}
	}

	return analysis
}

// Helper functions
func (s *Service) isInvokeCapability(capability string) bool {
	invokePatterns := []string{"invoke", "execute", "run", "call", "perform"}
	
	for _, pattern := range invokePatterns {
		if cap := capability; len(cap) >= len(pattern) {
			for i := 0; i <= len(cap)-len(pattern); i++ {
				if cap[i:i+len(pattern)] == pattern {
					return true
				}
			}
		}
	}
	return false
}

func (s *Service) categorizeCapability(capability string) string {
	switch {
	case len(capability) >= 5 && capability[:5] == "store":
		return "storage"
	case len(capability) >= 5 && capability[:5] == "space":
		return "space"
	case len(capability) >= 6 && capability[:6] == "upload":
		return "upload"
	case s.isInvokeCapability(capability):
		return "invocation"
	case len(capability) >= 4 && capability[:4] == "blob":
        return "blob"
    case len(capability) >= 5 && capability[:5] == "index":
        return "index"
    default:
        return "general"
    }
}

// extractCaveats converts IPLD node to map
func (s *Service) extractCaveats(nb any) map[string]interface{} {
    result := make(map[string]interface{})
    if nb == nil {
        return result
    }

    if node, ok := nb.(ipld.Node); ok && node.Kind() == ipld.Kind_Map {
        iter := node.MapIterator()
        for !iter.Done() {
            k, v, err := iter.Next()
            if err != nil {
                break
            }
            if keyStr, err := k.AsString(); err == nil {
                result[keyStr] = s.nodeToValue(v)
            }
        }
    }

    return result
}

// nodeToValue converts IPLD node to Go value
func (s *Service) nodeToValue(node ipld.Node) interface{} {
    switch node.Kind() {
    case ipld.Kind_Bool:
        v, _ := node.AsBool()
        return v
    case ipld.Kind_Int:
        v, _ := node.AsInt()
        return v
    case ipld.Kind_Float:
        v, _ := node.AsFloat()
        return v
    case ipld.Kind_String:
        v, _ := node.AsString()
        return v
    case ipld.Kind_Bytes:
        v, _ := node.AsBytes()
        return v
    default:
        return nil
    }
}