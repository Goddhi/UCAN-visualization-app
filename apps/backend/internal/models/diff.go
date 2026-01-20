package models

import "github.com/storacha/go-ucanto/ucan"

// DiffRequest comes from the frontend with two raw token strings
type DiffRequest struct {
	ParentToken string `json:"parent_token"`
	ChildToken  string `json:"child_token"`
}

// CapabilityDiff describes what happened to a single permission
type CapabilityDiff struct {
	// The capability as seen in the Child (Result)
	ChildCap map[string]interface{} `json:"child_cap"`

	// The capability in the Parent that authorized this (if found)
	ParentCap map[string]interface{} `json:"parent_cap,omitempty"`

	// Status: "UNCHANGED", "NARROWED", "ADDED" (Escalation), "REMOVED"
	Status string `json:"status"`

	Message string `json:"message"`
}

type DiffResponse struct {
	Diffs []CapabilityDiff `json:"diffs"`
}

// Helper to convert ucan.Capability to a JSON-friendly map
func CapToMap(c ucan.Capability[any]) map[string]interface{} {
	return map[string]interface{}{
		"can":  c.Can(),
		"with": c.With(),
		"nb":   c.Nb(),
	}
}