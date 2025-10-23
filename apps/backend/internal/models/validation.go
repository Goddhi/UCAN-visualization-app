package models

import "time"

// ValidationResult contains the complete validation outcome
type ValidationResult struct {
	Valid     bool              `json:"valid"`
	Chain     []ChainLink       `json:"chain"`
	RootCause *ValidationError  `json:"rootCause,omitempty"`
	Summary   ValidationSummary `json:"summary"`
}

// ValidationChainResult for validating complete chains
type ValidationChainResult struct {
	Valid        bool                   `json:"valid"`
	Chain        []ChainValidationLink  `json:"chain"`
	RootCause    *ValidationError       `json:"rootCause,omitempty"`
	Summary      ValidationSummary      `json:"summary"`
	ChainInfo    ChainInfo              `json:"chainInfo"`
}

// ChainValidationLink represents validation of a single link in a chain
type ChainValidationLink struct {
	Level         int                 `json:"level"`
	CID           string              `json:"cid"`
	Issuer        string              `json:"issuer"`
	Audience      string              `json:"audience"`
	Capabilities  []CapabilityInfo    `json:"capabilities"`
	Expiration    time.Time           `json:"expiration"`
	NotBefore     time.Time           `json:"notBefore"`
	Valid         bool                `json:"valid"`
	Issues        []ValidationIssue   `json:"issues,omitempty"`
	ProofValidation []ProofValidation `json:"proofValidation,omitempty"`
}

// ProofValidation represents validation of a specific proof
type ProofValidation struct {
	ProofCID     string            `json:"proofCid"`
	Valid        bool              `json:"valid"`
	Issues       []ValidationIssue `json:"issues,omitempty"`
	Capabilities []CapabilityInfo  `json:"capabilities"`
	Attenuation  AttenuationCheck  `json:"attenuation"`
}

// AttenuationCheck verifies that child capabilities are properly attenuated
type AttenuationCheck struct {
	Valid               bool     `json:"valid"`
	Issues              []string `json:"issues,omitempty"`
	ResourceMatch       bool     `json:"resourceMatch"`
	AbilityMatch        bool     `json:"abilityMatch"`
	CaveatProperlyAdded bool     `json:"caveatProperlyAdded"`
}

// ChainLink represents a single link in the validation chain
type ChainLink struct {
	Level      int              `json:"level"`
	CID        string           `json:"cid"`
	Issuer     string           `json:"issuer"`
	Audience   string           `json:"audience"`
	Capability CapabilityInfo   `json:"capability"`
	Expiration time.Time        `json:"expiration"`
	NotBefore  time.Time        `json:"notBefore"`
	Valid      bool             `json:"valid"`
	Issues     []ValidationIssue `json:"issues,omitempty"`
}

// ValidationIssue represents a specific validation problem
type ValidationIssue struct {
	Type     string                 `json:"type"`
	Message  string                 `json:"message"`
	Severity string                 `json:"severity"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// ValidationError represents the root cause of validation failure
type ValidationError struct {
	Type    string     `json:"type"`
	Message string     `json:"message"`
	Link    *LinkInfo  `json:"link,omitempty"`
}

// LinkInfo contains minimal link information
type LinkInfo struct {
	Issuer   string `json:"issuer"`
	Audience string `json:"audience"`
}

// ValidationSummary provides statistics about the validation
type ValidationSummary struct {
	TotalLinks   int `json:"totalLinks"`
	ValidLinks   int `json:"validLinks"`
	InvalidLinks int `json:"invalidLinks"`
	WarningCount int `json:"warningCount"`
}