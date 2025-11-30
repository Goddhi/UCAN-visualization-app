package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// RootHandler handles the root endpoint
func RootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "UCAN Visualizer API",
		"version":   "1.0.0",
		"status":    "running",
		"timestamp": time.Now().Format(time.RFC3339),
		"endpoints": map[string]interface{}{
			"health": "GET /health",
			"parse": map[string]string{
				"delegation":      "POST /api/parse/delegation",
				"delegation_file": "POST /api/parse/delegation/file",
				"chain":           "POST /api/parse/chain",
				"chain_file":      "POST /api/parse/chain/file",
				"invocation":      "POST /api/parse/invocation",
				"invocation_file": "POST /api/parse/invocation/file",
			},
			"validate": map[string]string{
				"chain":      "POST /api/validate/chain",
				"chain_file": "POST /api/validate/chain/file",
			},
			"graph": map[string]string{
				"delegation":      "POST /api/graph/delegation",
				"delegation_file": "POST /api/graph/delegation/file",
				"invocation":      "POST /api/graph/invocation",
				"invocation_file": "POST /api/graph/invocation/file",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}