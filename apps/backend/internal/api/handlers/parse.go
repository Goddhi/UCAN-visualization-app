package handlers

import (
	"encoding/json"
	"fmt"
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

// ParseDelegation handles POST /api/parse/delegation
func (h *ParseHandler) ParseDelegation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse delegation request from %s", r.RemoteAddr)

	tokenBytes, err := h.extractTokenFromRequest(r)
	if err != nil {
		log.Printf("[ERROR] Failed to extract token: %v", err)
		respondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("[DEBUG] Processing token of length %d bytes", len(tokenBytes))

	result, err := h.parser.ParseDelegation(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse delegation", err)
		return
	}

	log.Printf("[INFO] Successfully parsed delegation: %s", result.CID)
	respondJSON(w, http.StatusOK, result)
}

// ParseChain handles POST /api/parse/chain
func (h *ParseHandler) ParseChain(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse chain request from %s", r.RemoteAddr)

	tokenBytes, err := h.extractTokenFromRequest(r)
	if err != nil {
		log.Printf("[ERROR] Failed to extract token: %v", err)
		respondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("[DEBUG] Processing chain token of length %d bytes", len(tokenBytes))

	result, err := h.parser.ParseDelegationChain(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Chain parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse delegation chain", err)
		return
	}

	log.Printf("[INFO] Successfully parsed delegation chain: %d delegations", len(result))
	respondJSON(w, http.StatusOK, result)
}

// ParseInvocation handles POST /api/parse/invocation
func (h *ParseHandler) ParseInvocation(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse invocation request from %s", r.RemoteAddr)

	tokenBytes, err := h.extractTokenFromRequest(r)
	if err != nil {
		log.Printf("[ERROR] Failed to extract token: %v", err)
		respondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("[DEBUG] Processing invocation token of length %d bytes", len(tokenBytes))

	result, err := h.parser.ParseInvocation(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Invocation parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse invocation", err)
		return
	}

	log.Printf("[INFO] Successfully parsed invocation: is_invocation=%v, task_type=%s", 
		result.IsInvocation, 
		func() string {
			if result.Task != nil {
				return result.Task.TaskType
			}
			return "none"
		}())
	respondJSON(w, http.StatusOK, result)
}

// ParseFile handles POST /api/parse/delegation/file
func (h *ParseHandler) ParseFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse file upload request from %s", r.RemoteAddr)

	tokenBytes, err := h.extractTokenFromFile(r)
	if err != nil {
		log.Printf("[ERROR] Failed to extract token from file: %v", err)
		respondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("[DEBUG] Processing file token of length %d bytes", len(tokenBytes))

	result, err := h.parser.ParseDelegation(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse delegation", err)
		return
	}

	log.Printf("[INFO] Successfully parsed delegation from file: %s", result.CID)
	respondJSON(w, http.StatusOK, result)
}

// ParseChainFile handles POST /api/parse/chain/file
func (h *ParseHandler) ParseChainFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse chain file upload request from %s", r.RemoteAddr)

	tokenBytes, err := h.extractTokenFromFile(r)
	if err != nil {
		log.Printf("[ERROR] Failed to extract token from file: %v", err)
		respondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("[DEBUG] Processing chain file token of length %d bytes", len(tokenBytes))

	result, err := h.parser.ParseDelegationChain(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Chain parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse delegation chain", err)
		return
	}

	log.Printf("[INFO] Successfully parsed delegation chain from file: %d delegations", len(result))
	respondJSON(w, http.StatusOK, result)
}

// ParseInvocationFile handles POST /api/parse/invocation/file
func (h *ParseHandler) ParseInvocationFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Parse invocation file upload request from %s", r.RemoteAddr)

	tokenBytes, err := h.extractTokenFromFile(r)
	if err != nil {
		log.Printf("[ERROR] Failed to extract token from file: %v", err)
		respondError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	log.Printf("[DEBUG] Processing invocation file token of length %d bytes", len(tokenBytes))

	result, err := h.parser.ParseInvocation(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Invocation parse failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to parse invocation", err)
		return
	}

	log.Printf("[INFO] Successfully parsed invocation from file: is_invocation=%v", result.IsInvocation)
	respondJSON(w, http.StatusOK, result)
}

// extractTokenFromRequest extracts token from JSON request body
func (h *ParseHandler) extractTokenFromRequest(r *http.Request) ([]byte, error) {
	var req models.ParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	if req.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Normalize token - primarily for JWT format
	tokenBytes, err := utils.NormalizeToken(req.Token, req.Format)
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	return tokenBytes, nil
}

// extractTokenFromFile extracts token from uploaded file
func (h *ParseHandler) extractTokenFromFile(r *http.Request) ([]byte, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("file is required: %w", err)
	}
	defer file.Close()

	tokenBytes, err := utils.ReadUploadedFile(file, header)
	if err != nil {
		return nil, fmt.Errorf("failed to read uploaded file: %w", err)
	}

	// Basic validation for UCAN token files
	if err := utils.ValidateTokenFormat(string(tokenBytes)); err != nil {
		return nil, fmt.Errorf("invalid UCAN token file: %w", err)
	}

	return tokenBytes, nil
}