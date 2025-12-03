package models

import "time"

// Core delegation models
type DelegationResponse struct {
	Issuer       string              `json:"issuer"`
	Audience     string              `json:"audience"`
	Capabilities []CapabilityInfo    `json:"capabilities"`
	Proofs       []ProofInfo         `json:"proofs"`
	Expiration   time.Time           `json:"expiration,omitempty"`
	NotBefore    time.Time           `json:"notBefore,omitempty"`
	IssuedAt     time.Time           `json:"issuedAt,omitempty"`
	Facts        []interface{}       `json:"facts,omitempty"`
	Nonce        string              `json:"nonce,omitempty"`
	Signature    SignatureInfo       `json:"signature"`
	CID          string              `json:"cid"`
	Level        int                 `json:"level"`
}

// Enhanced capability model
type CapabilityInfo struct {
	With     string                 `json:"with"`
	Can      string                 `json:"can"`
	Nb       map[string]interface{} `json:"nb"`
	Category string                 `json:"category"` // storage, space, upload, invocation, etc.
}

// Enhanced proof model
type ProofInfo struct {
	CID   string `json:"cid"`
	Index int    `json:"index"`
	Type  string `json:"type"` // delegation, invocation, receipt
}

type SignatureInfo struct {
	Algorithm string `json:"algorithm"`
	Verified  bool   `json:"verified"`        
	Valid     bool   `json:"valid"`           
	Error     string `json:"error,omitempty"`
}
// Enhanced invocation models
type InvocationResponse struct {
	Delegation         *DelegationResponse  `json:"delegation"`
	IsInvocation       bool                 `json:"isInvocation"`
	Task               *TaskInfo            `json:"task,omitempty"`
	InvocationAnalysis *InvocationAnalysis  `json:"invocationAnalysis"`
	CapabilityAnalysis *CapabilityAnalysis  `json:"capabilityAnalysis"`
}

type TaskInfo struct {
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	Constraints map[string]interface{} `json:"constraints"`
	Issuer      string                 `json:"issuer"`
	Target      string                 `json:"target"`
	TaskType    string                 `json:"taskType"` // invocation, delegation
	Permissions []string               `json:"permissions"`
}

// Comprehensive invocation analysis
type InvocationAnalysis struct {
	IsInvocation        bool                   `json:"isInvocation"`
	HasInvokeCapability bool                   `json:"hasInvokeCapability"`
	TaskType           string                 `json:"taskType"` // invocation, delegation
	PrimaryAction      string                 `json:"primaryAction"`
	TargetResource     string                 `json:"targetResource"`
	InvokePatterns     []string               `json:"invokePatterns"`
	RequiredPermissions []string               `json:"requiredPermissions"`
	Constraints        map[string]interface{} `json:"constraints"`
}

// Comprehensive capability analysis
type CapabilityAnalysis struct {
	Categories    map[string][]CapabilityInfo `json:"categories"`
	TotalCount    int                         `json:"totalCount"`
	InvokeCount   int                         `json:"invokeCount"`
	DelegateCount int                         `json:"delegateCount"`
	Permissions   []string                    `json:"permissions"`
	Resources     []string                    `json:"resources"`
}


// Chain analysis models
type ChainInfo struct {
	TotalLevels int               `json:"totalLevels"`
	IsComplete  bool              `json:"isComplete"`
	RootCID     string            `json:"rootCid"`
	LeafCIDs    []string          `json:"leafCids"`
	Principals  []PrincipalInfo   `json:"principals"`
	Timeline    []TimelineEvent   `json:"timeline"`
	ProofChain  *ProofChain       `json:"proofChain,omitempty"`
}

type ProofChain struct {
	Root       string       `json:"root"`
	Levels     []ProofLevel `json:"levels"`
	IsComplete bool         `json:"isComplete"`
	TotalDepth int          `json:"totalDepth"`
}

type ProofLevel struct {
	Level       int                    `json:"level"`
	Delegations []*DelegationResponse  `json:"delegations"`
	ProofLinks  []string               `json:"proofLinks"`
}

// Validation models
type ValidationResponse struct {
	Valid     bool              `json:"valid"`
	Chain     []DelegationInfo  `json:"chain"`
	Errors    []ValidationError `json:"errors,omitempty"`
	RootCause *ValidationError  `json:"rootCause,omitempty"`
	Summary   ValidationSummary `json:"summary"`
}

type DelegationInfo struct {
	CID      string `json:"cid"`
	Level    int    `json:"level"`
	Valid    bool   `json:"valid"`
	Issuer   string `json:"issuer"`
	Audience string `json:"audience"`
}

