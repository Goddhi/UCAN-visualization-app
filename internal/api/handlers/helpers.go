package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/goddhi/ucan-visualizer/internal/models"
)

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string, err error) {
	errResp := models.ErrorResponse{
		Error:     http.StatusText(status),
		Message:   message,
		Timestamp: time.Now(),
	}

	if err != nil {
		errResp.Details = map[string]interface{}{
			"error": err.Error(),
		}
	}

	respondJSON(w, status, errResp)
}

// HealthCheck returns the health status of the service
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"time":    time.Now(),
		"service": "ucan-visualizer",
		"version": "1.0.0",
	})
}

// isValidUCANFile checks if the file has a valid extension for UCAN tokens
func isValidUCANFile(filename string) bool {
	if filename == "" {
		return true // Allow files without names
	}

	ext := strings.ToLower(filepath.Ext(filename))
	
	// Valid extensions for UCAN CAR files
	validExtensions := []string{
		".car",    
		".ucan",    
		".cbor",    
		"",        
	}

	for _, valid := range validExtensions {
		if ext == valid {
			return true
		}
	}

	return false
}