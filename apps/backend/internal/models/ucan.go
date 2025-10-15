package models

import "time"

type DelegationResponse struct {
	Issuer       string           `json:"issuer"`
	Audience     string           `json:"audience"`
	Subject      string           `json:"subject,omitempty"`
	Capabilities []CapabilityInfo `json:"capabilities"`
	Proofs       []ProofInfo      `json:"proofs"`
	Expiration   time.Time        `json:"expiration"`
	NotBefore    time.Time        `json:"notBefore"`
	Facts        []interface{}    `json:"facts,omitempty"`
	Nonce        string           `json:"nonce,omitempty"`
	Signature    SignatureInfo    `json:"signature"`
	CID          string           `json:"cid"`
}

// CapabilityInfo represents a UCAN capability
type CapabilityInfo struct {
	With string                 `json:"with"` // Resource URI
	Can  string                 `json:"can"`  // Ability
	Nb   map[string]interface{} `json:"nb"`   // Caveats
}

// ProofInfo represents a proof in the chain
type ProofInfo struct {
	CID          string           `json:"cid"`
	Issuer       string           `json:"issuer"`
	Audience     string           `json:"audience"`
	Capabilities []CapabilityInfo `json:"capabilities"`
	Expiration   time.Time        `json:"expiration"`
	NotBefore    time.Time        `json:"notBefore"`
	Proofs       []ProofInfo      `json:"proofs,omitempty"` // Recursive
}

// SignatureInfo represents signature metadata
type SignatureInfo struct {
	Algorithm string `json:"algorithm"`
	Valid     *bool  `json:"valid,omitempty"`
}