package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/validator"
	"github.com/goddhi/ucan-visualizer/pkg/utils"
)

type ValidateHandler struct {
	validator *validator.Service
}

func NewValidateHandler() *ValidateHandler {
	return &ValidateHandler{
		validator: validator.NewService(),
	}
}

// ValidateChain handles JSON requests with token string
func (h *ValidateHandler) ValidateChain(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Validate request from %s", r.RemoteAddr)

	var req models.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[ERROR] Failed to decode request: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Token == "" {
		log.Printf("[WARN] Empty token in request")
		respondError(w, http.StatusBadRequest, "Token is required", nil)
		return
	}

	// Auto-detect format and normalize to bytes
	tokenBytes, err := utils.NormalizeToken(req.Token, req.Format)
	if err != nil {
		log.Printf("[ERROR] Failed to normalize token: %v", err)
		respondError(w, http.StatusBadRequest, 
			"Invalid token format. Supported formats: base64, hex. "+
			"Or upload a .car file.", err)
		return
	}

	log.Printf("[DEBUG] Validating token of length %d bytes", len(tokenBytes))
	result, err := h.validator.ValidateChain(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Validation failed", err)
		return
	}

	log.Printf("[INFO] Validation complete: Valid=%v", result.Valid)
	respondJSON(w, http.StatusOK, result)
}

// ValidateFile handles multipart file uploads
func (h *ValidateHandler) ValidateFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Validate file upload request from %s", r.RemoteAddr)

	// Parse multipart form (10MB max)
	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		log.Printf("[ERROR] Failed to parse multipart form: %v", err)
		respondError(w, http.StatusBadRequest, "Failed to parse multipart form", err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[ERROR] No file provided: %v", err)
		respondError(w, http.StatusBadRequest, "No file provided. Use 'file' as the form field name.", err)
		return
	}
	defer file.Close()

	log.Printf("[DEBUG] Received file: %s (%d bytes)", header.Filename, header.Size)
	// Read file contents
	tokenBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[ERROR] Failed to read file: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to read file", err)
		return
	}

	log.Printf("[DEBUG] Read %d bytes from file", len(tokenBytes))

	// Validate the delegation
	result, err := h.validator.ValidateChain(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Validation failed", err)
		return
	}

	log.Printf("[INFO] Validation complete from file: Valid=%v", result.Valid)
	respondJSON(w, http.StatusOK, result)
}