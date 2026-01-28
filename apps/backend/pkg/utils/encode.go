package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// UCANClaims represents the standard UCAN payload
type UCANClaims struct {
	Issuer    string                   `json:"iss"`
	Audience  string                   `json:"aud"`
	Expiry    int64                    `json:"exp"`
	NotBefore int64                    `json:"nbf"`
	Nonce     string                   `json:"nnc"`
	Facts     []interface{}            `json:"fct"`
	Proofs    []string                 `json:"prf"`
	Att       []map[string]interface{} `json:"att"` // Capabilities
	Cid       string                   `json:"cid,omitempty"`
}

// ParsedJWT holds the raw data we extracted
type ParsedJWT struct {
	Header    map[string]interface{}
	Claims    UCANClaims
	Signature []byte
}

// ParseUnverifiedJWT decodes a standard JWT string (ey...)
func ParseUnverifiedJWT(tokenString string) (*ParsedJWT, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	// 1. Decode Header
	headerBytes, err := parseBase64(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}
	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal header: %w", err)
	}

	// 2. Decode Payload (Claims)
	payloadBytes, err := parseBase64(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}
	var claims UCANClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// 3. Decode Signature
	sigBytes, err := parseBase64(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	return &ParsedJWT{
		Header:    header,
		Claims:    claims,
		Signature: sigBytes,
	}, nil
}

// ParseUnverifiedCBOR decodes a Raw UCAN Block (glhA...)
func ParseUnverifiedCBOR(data []byte) (*ParsedJWT, error) {
	// 1. Decode DAG-CBOR into a Generic Node
	nb := basicnode.Prototype.Any.NewBuilder()
	if err := dagcbor.Decode(nb, bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("failed to decode DAG-CBOR: %w", err)
	}
	node := nb.Build()

	if node.Kind() != ipld.Kind_List {
		return nil, fmt.Errorf("invalid UCAN block: expected array, got %v", node.Kind())
	}

	claims := UCANClaims{}
	var sigBytes []byte
	var header map[string]interface{}

	iter := node.ListIterator()
	for !iter.Done() {
		_, item, _ := iter.Next()

		if item.Kind() == ipld.Kind_Map {
			// Check immediate keys for UCAN Version Wrapper
			mIter := item.MapIterator()
			foundNested := false
			
			for !mIter.Done() {
				k, v, _ := mIter.Next()
				kStr, _ := k.AsString()

				if strings.HasPrefix(kStr, "ucan/") {
					log.Printf("[DEBUG] Found UCAN Version Key: %s", kStr)
					claims = extractClaims(v)
					foundNested = true
				}
			}

			if !foundNested && claims.Issuer == "" {
				tempClaims := extractClaims(item)
				if tempClaims.Issuer != "" || tempClaims.Audience != "" {
					claims = tempClaims
				} else {
					header = extractMap(item)
				}
			}

		} else if item.Kind() == ipld.Kind_Bytes {
			b, _ := item.AsBytes()
			if len(b) == 64 {
				sigBytes = b
			}
		}
	}

	if claims.Issuer == "" && claims.Audience == "" {
		return nil, fmt.Errorf("valid CBOR array found, but contained no recognizable UCAN Payload")
	}

	return &ParsedJWT{
		Header:    header,
		Claims:    claims,
		Signature: sigBytes,
	}, nil
}

// extractClaims converts an IPLD Node (Map) into our UCANClaims struct
func extractClaims(node ipld.Node) UCANClaims {
	claims := UCANClaims{}
	if node.Kind() != ipld.Kind_Map {
		return claims
	}
	
	// Variables to hold iso-ucan specific fields
	var cmd string
	var sub string
	var pol interface{}

	pIter := node.MapIterator()
	for !pIter.Done() {
		k, v, _ := pIter.Next()
		keyStr, _ := k.AsString()

		log.Printf("[DEBUG] Claims Key Found: %s", keyStr)

		switch keyStr {
		case "iss": claims.Issuer, _ = v.AsString()
		case "aud": claims.Audience, _ = v.AsString()
		case "exp": exp, _ := v.AsInt(); claims.Expiry = exp
		case "nbf": nbf, _ := v.AsInt(); claims.NotBefore = nbf
		case "nnc": claims.Nonce, _ = v.AsString()
		
		// --- CAPTURE ISO-UCAN FIELDS ---
		case "cmd":
			cmd, _ = v.AsString()
		case "sub":
			sub, _ = v.AsString()
		case "pol":
			pol = nodeToValue(v)
			// Keep adding to facts as well, as it is technically a fact/assertion
			claims.Facts = append(claims.Facts, map[string]interface{}{"pol": pol})

		case "prf":
			if v.Kind() == ipld.Kind_List {
				lIter := v.ListIterator()
				for !lIter.Done() {
					_, pNode, _ := lIter.Next()
					if link, err := pNode.AsLink(); err == nil {
						claims.Proofs = append(claims.Proofs, link.String())
					} else if str, err := pNode.AsString(); err == nil {
						claims.Proofs = append(claims.Proofs, str)
					}
				}
			}

		case "att", "capabilities", "caps": 
			if v.Kind() == ipld.Kind_List {
				lIter := v.ListIterator()
				for !lIter.Done() {
					_, capNode, _ := lIter.Next()
					if capNode.Kind() == ipld.Kind_Map {
						capMap := make(map[string]interface{})
						cIter := capNode.MapIterator()
						for !cIter.Done() {
							ck, cv, _ := cIter.Next()
							ckStr, _ := ck.AsString()
							capMap[ckStr] = nodeToValue(cv)
						}
						claims.Att = append(claims.Att, capMap)
					}
				}
			}
		}
	}

	// --- SYNTHESIZE CAPABILITY ---
	// If we found 'cmd' (Action) but no standard 'att' array, map it!
	if len(claims.Att) == 0 && cmd != "" {
		log.Printf("[DEBUG] Mapping iso-ucan 'cmd' to standard Capability...")
		capMap := map[string]interface{}{
			"can": cmd,  // cmd -> can
			"with": sub, // sub -> with
		}
		// If there is a policy, wrap it as a caveat
		if pol != nil {
			capMap["nb"] = map[string]interface{}{"pol": pol}
		}
		claims.Att = append(claims.Att, capMap)
	}

	return claims
}

func extractMap(node ipld.Node) map[string]interface{} {
	m := make(map[string]interface{})
	iter := node.MapIterator()
	for !iter.Done() {
		k, v, _ := iter.Next()
		kStr, _ := k.AsString()
		m[kStr] = nodeToValue(v)
	}
	return m
}

func parseBase64(input string) ([]byte, error) {
	if l := len(input) % 4; l > 0 {
		input += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(input)
}

func nodeToValue(node ipld.Node) interface{} {
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
	case ipld.Kind_Link:
		v, _ := node.AsLink()
		return v.String()
	case ipld.Kind_Map:
		return extractMap(node)
	case ipld.Kind_List:
		var l []interface{}
		iter := node.ListIterator()
		for !iter.Done() {
			_, v, _ := iter.Next()
			l = append(l, nodeToValue(v))
		}
		return l
	default:
		return nil
	}
}

