package parser

import (
	"fmt"
	"time"

	"github.com/ipld/go-ipld-prime"
	"github.com/storacha/go-ucanto/core/dag/blockstore"
	"github.com/storacha/go-ucanto/core/delegation"
	"github.com/goddhi/ucan-visualizer/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// ParseDelegation parses a UCAN delegation token (CAR format)
func (s *Service) ParseDelegation(tokenBytes []byte) (*models.DelegationResponse, error) {
	// Use delegation.Extract to parse CAR-encoded delegation
	del, err := delegation.Extract(tokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract delegation: %w", err)
	}

	issuer := del.Issuer().DID().String()
	audience := del.Audience().DID().String()

	var subject string

	// Parse capabilities
	capabilities := []models.CapabilityInfo{}
	for _, cap := range del.Capabilities() {
		capInfo := models.CapabilityInfo{
			With: cap.With(),
			Can:  cap.Can(),
			Nb:   s.extractCaveats(cap.Nb()),
		}
		capabilities = append(capabilities, capInfo)
	}

	// Parse proofs, these are Links, not full delegations
	proofs := []models.ProofInfo{}
	for _, proofLink := range del.Proofs() {
		proofInfo := models.ProofInfo{
			CID: proofLink.String(),
		}
		proofs = append(proofs, proofInfo)
	}

	// Extract time bounds
	exp := del.Expiration()
	var expiration time.Time
	if exp != nil {
		expiration = time.Unix(int64(*exp), 0)
	}

	nbf := del.NotBefore()
	notBefore := time.Unix(int64(nbf), 0)

	nonce := string(del.Nonce())

	facts := []interface{}{}
	for _, fact := range del.Facts() {
		facts = append(facts, fact)
	}

	return &models.DelegationResponse{
		Issuer:       issuer,
		Audience:     audience,
		Subject:      subject,
		Capabilities: capabilities,
		Proofs:       proofs,
		Expiration:   expiration,
		NotBefore:    notBefore,
		Facts:        facts,
		Nonce:        nonce,
		Signature: models.SignatureInfo{
			Algorithm: "EdDSA",
			Valid:     nil,
		},
		CID: del.Link().String(),
	}, nil
}

// ParseDelegationWithProofs parses a delegation and recursively resolves proofs from the blockstore
func (s *Service) ParseDelegationWithProofs(tokenBytes []byte) (*models.DelegationResponse, error) {
	del, err := delegation.Extract(tokenBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to extract delegation: %w", err)
	}

	// Extract issuer and audience
	issuer := del.Issuer().DID().String()
	audience := del.Audience().DID().String()

	// Extract subject (not available at delegation level in go-ucanto)
	var subject string

	// Parse capabilities
	capabilities := []models.CapabilityInfo{}
	for _, cap := range del.Capabilities() {
		capInfo := models.CapabilityInfo{
			With: cap.With(),
			Can:  cap.Can(),
			Nb:   s.extractCaveats(cap.Nb()),
		}
		capabilities = append(capabilities, capInfo)
	}

	// Create a blockstore from the delegation's blocks to resolve proofs
	br, err := blockstore.NewBlockReader(
		blockstore.WithBlocksIterator(del.Blocks()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating block reader: %w", err)
	}

	// Resolve each proof
	proofs := []models.ProofInfo{}
	for _, proofLink := range del.Proofs() {
		proofDel, err := delegation.NewDelegationView(proofLink, br)
		// if not resolvable, return proof CID link
		if err != nil {
			proofs = append(proofs, models.ProofInfo{
				CID:      proofLink.String(),
				Issuer:   "unresolved",
				Audience: "unresolved",
			})
			continue
		}

		// if successfully resolved the proof, Extract its data
		proofExp := proofDel.Expiration()
		var proofExpiration time.Time
		if proofExp != nil {
			proofExpiration = time.Unix(int64(*proofExp), 0)
		}

		proofNbf := proofDel.NotBefore()
		proofNotBefore := time.Unix(int64(proofNbf), 0)

		proofInfo := models.ProofInfo{
			CID:          proofLink.String(),
			Issuer:       proofDel.Issuer().DID().String(),
			Audience:     proofDel.Audience().DID().String(),
			Capabilities: []models.CapabilityInfo{},
			Expiration:   proofExpiration,
			NotBefore:    proofNotBefore,
		}

		// Extract capabilities from the proof
		for _, cap := range proofDel.Capabilities() {
			proofInfo.Capabilities = append(proofInfo.Capabilities, models.CapabilityInfo{
				With: cap.With(),
				Can:  cap.Can(),
				Nb:   s.extractCaveats(cap.Nb()),
			})
		}

		// Recursively resolve nested proofs if they exist
		if len(proofDel.Proofs()) > 0 {
			nestedProofs := []models.ProofInfo{}
			for _, nestedLink := range proofDel.Proofs() {
				nestedDel, err := delegation.NewDelegationView(nestedLink, br)
				if err != nil {
					// Nested proof not found
					nestedProofs = append(nestedProofs, models.ProofInfo{
						CID:      nestedLink.String(),
						Issuer:   "unresolved",
						Audience: "unresolved",
					})
					continue
				}

				nestedExp := nestedDel.Expiration()
				var nestedExpiration time.Time
				if nestedExp != nil {
					nestedExpiration = time.Unix(int64(*nestedExp), 0)
				}

				nestedNbf := nestedDel.NotBefore()
				nestedNotBefore := time.Unix(int64(nestedNbf), 0)

				nestedInfo := models.ProofInfo{
					CID:          nestedLink.String(),
					Issuer:       nestedDel.Issuer().DID().String(),
					Audience:     nestedDel.Audience().DID().String(),
					Capabilities: []models.CapabilityInfo{},
					Expiration:   nestedExpiration,
					NotBefore:    nestedNotBefore,
				}

				for _, cap := range nestedDel.Capabilities() {
					nestedInfo.Capabilities = append(nestedInfo.Capabilities, models.CapabilityInfo{
						With: cap.With(),
						Can:  cap.Can(),
						Nb:   s.extractCaveats(cap.Nb()),
					})
				}

				nestedProofs = append(nestedProofs, nestedInfo)
			}
			proofInfo.Proofs = nestedProofs
		}

		proofs = append(proofs, proofInfo)
	}

	// Extract time bounds
	exp := del.Expiration()
	var expiration time.Time
	if exp != nil {
		expiration = time.Unix(int64(*exp), 0)
	}

	nbf := del.NotBefore()
	notBefore := time.Unix(int64(nbf), 0)

	// Extract nonce
	nonce := string(del.Nonce())

	// Extract facts
	facts := []interface{}{}
	for _, fact := range del.Facts() {
		facts = append(facts, fact)
	}

	return &models.DelegationResponse{
		Issuer:       issuer,
		Audience:     audience,
		Subject:      subject,
		Capabilities: capabilities,
		Proofs:       proofs,
		Expiration:   expiration,
		NotBefore:    notBefore,
		Facts:        facts,
		Nonce:        nonce,
		Signature: models.SignatureInfo{
			Algorithm: "EdDSA",
			Valid:     nil,
		},
		CID: del.Link().String(),
	}, nil
}

// extractCaveats extracts caveats from the capability's Nb field
func (s *Service) extractCaveats(nb any) map[string]interface{} {
	result := make(map[string]interface{})

	if nb == nil {
		return result
	}

	if node, ok := nb.(ipld.Node); ok {
		if node.Kind() == ipld.Kind_Map {
			iter := node.MapIterator()
			for !iter.Done() {
				k, v, err := iter.Next()
				if err != nil {
					break
				}

				keyStr, err := k.AsString()
				if err != nil {
					continue
				}

				result[keyStr] = s.nodeToValue(v)
			}
		}
	}

	return result
}

// nodeToValue converts IPLD node to Go value
func (s *Service) nodeToValue(node ipld.Node) interface{} {
	if node == nil {
		return nil
	}

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
	case ipld.Kind_List:
		var result []interface{}
		iter := node.ListIterator()
		for !iter.Done() {
			_, val, err := iter.Next()
			if err != nil {
				break
			}
			result = append(result, s.nodeToValue(val))
		}
		return result
	case ipld.Kind_Map:
		result := make(map[string]interface{})
		iter := node.MapIterator()
		for !iter.Done() {
			k, v, err := iter.Next()
			if err != nil {
				break
			}
			keyStr, _ := k.AsString()
			result[keyStr] = s.nodeToValue(v)
		}
		return result
	default:
		return nil
	}
}