package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":   "ucan-visualizer",
		"status":    "healthy",
		"time":      time.Now(),
		"version":   "1.0.0",
	}
	respondJSON(w, http.StatusOK, response)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string, err error) {
	response := map[string]interface{}{
		"error":     http.StatusText(status),
		"message":   message,
		"timestamp": time.Now(),
	}
	
	if err != nil {
		response["details"] = map[string]interface{}{
			"error": err.Error(),
		}
	}
	
	respondJSON(w, status, response)
}