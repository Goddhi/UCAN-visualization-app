package models

import "time"

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error     string                 `json:"error"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	RequestID string                 `json:"requestId,omitempty"`
}
type ParseRequest struct {

	Token string `json:"token"`
	
	Format string `json:"format,omitempty"`
}

type ValidateRequest struct {
	Token  string `json:"token"`
	Format string `json:"format,omitempty"`
}

type GraphRequest struct {
	Token  string `json:"token"`
	Format string `json:"format,omitempty"`
}