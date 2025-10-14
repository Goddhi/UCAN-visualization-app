package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/parser"
	"github.com/goddhi/ucan-visualizer/pkg/utils"
)

type ParseHandler struct {
	parser *parser.Service
}

func NewParseHandler() *ParseHandler {
	return &ParseHandler{
		parser: parser.NewService(),
	}
}

// ParseDelegation handles JSON requests with token string
func (h *ParseHandler) ParseDelegation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse request from %s", r.RemoteAddr)

	var req models.ParseRequest
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

	log.Printf("[DEBUG] Parsing token of length %d bytes", len(tokenBytes))
	
	result, err := h.parser.ParseDelegation(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse delegation", err)
		return
	}

	log.Printf("[INFO] Successfully parsed delegation: %s", result.CID)
	respondJSON(w, http.StatusOK, result)
}

// ParseFile handles multipart file uploads
func (h *ParseHandler) ParseFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse file upload request from %s", r.RemoteAddr)

	// Parse multipart form (10MB max)
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("[ERROR] Failed to parse multipart form: %v", err)
		respondError(w, http.StatusBadRequest, "Failed to parse multipart form", err)
		return
	}

	// Get the file from the form
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

	// Parse the delegation
	result, err := h.parser.ParseDelegation(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse delegation", err)
		return
	}

	log.Printf("[INFO] Successfully parsed delegation from file: %s", result.CID)
	respondJSON(w, http.StatusOK, result)
}