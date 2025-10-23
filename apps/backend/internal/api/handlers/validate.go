package handlers

import (
	"encoding/json"
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

// ValidateChain handles POST /api/validate/chain
func (h *ValidateHandler) ValidateChain(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Validate chain request from %s", r.RemoteAddr)

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

	tokenBytes, err := utils.NormalizeToken(req.Token, req.Format)
	if err != nil {
		log.Printf("[ERROR] Failed to normalize token: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid token format", err)
		return
	}

	log.Printf("[DEBUG] Validating chain for token of length %d bytes", len(tokenBytes))

	result, err := h.validator.ValidateChain(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Validation failed", err)
		return
	}

	log.Printf("[INFO] Successfully validated chain: valid=%v, total_links=%d, valid_links=%d", 
		result.Valid, result.Summary.TotalLinks, result.Summary.ValidLinks)
	respondJSON(w, http.StatusOK, result)
}

// ValidateFile handles POST /api/validate/chain/file
func (h *ValidateHandler) ValidateFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Validate file upload request from %s", r.RemoteAddr)

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[ERROR] Failed to get file from form: %v", err)
		respondError(w, http.StatusBadRequest, "File is required", err)
		return
	}
	defer file.Close()

	tokenBytes, err := utils.ReadUploadedFile(file, header)
	if err != nil {
		log.Printf("[ERROR] Failed to read uploaded file: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid file", err)
		return
	}

	// Validate file content
	if err := utils.IsValidUCANFile(tokenBytes, header.Filename); err != nil {
		log.Printf("[ERROR] Invalid UCAN file: %v", err)
		respondError(w, http.StatusBadRequest, "Invalid UCAN file", err)
		return
	}

	log.Printf("[DEBUG] Validating chain for file %s (%d bytes)", 
		header.Filename, len(tokenBytes))

	result, err := h.validator.ValidateChain(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Validation failed: %v", err)
		respondError(w, http.StatusInternalServerError, "Validation failed", err)
		return
	}

	log.Printf("[INFO] Successfully validated chain from file: valid=%v, total_links=%d, valid_links=%d", 
		result.Valid, result.Summary.TotalLinks, result.Summary.ValidLinks)
	respondJSON(w, http.StatusOK, result)
}