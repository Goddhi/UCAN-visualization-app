package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/graph"
	"github.com/goddhi/ucan-visualizer/pkg/utils"
)

type GraphHandler struct {
	graph *graph.Service
}

func NewGraphHandler() *GraphHandler {
	return &GraphHandler{
		graph: graph.NewService(),
	}
}

// GenerateGraph handles POST /api/graph/delegation
func (h *GraphHandler) GenerateGraph(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Graph generation request from %s", r.RemoteAddr)

	var req models.GraphRequest
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

	log.Printf("[DEBUG] Generating delegation graph for token of length %d bytes", len(tokenBytes))
	
	result, err := h.graph.GenerateDelegationGraph(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Graph generation failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate graph", err)
		return
	}

	log.Printf("[INFO] Successfully generated delegation graph: %d nodes, %d edges", 
		len(result.Nodes), len(result.Edges))
	respondJSON(w, http.StatusOK, result)
}

// GenerateGraphFile handles POST /api/graph/delegation/file
func (h *GraphHandler) GenerateGraphFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Graph file generation request from %s", r.RemoteAddr)

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

	log.Printf("[DEBUG] Generating delegation graph for file %s (%d bytes)", 
		header.Filename, len(tokenBytes))
	
	result, err := h.graph.GenerateDelegationGraph(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Graph generation failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate graph", err)
		return
	}

	log.Printf("[INFO] Successfully generated delegation graph from file: %d nodes, %d edges", 
		len(result.Nodes), len(result.Edges))
	respondJSON(w, http.StatusOK, result)
}

// GenerateInvocationGraph handles POST /api/graph/invocation
func (h *GraphHandler) GenerateInvocationGraph(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Invocation graph generation request from %s", r.RemoteAddr)

	var req models.GraphRequest
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

	log.Printf("[DEBUG] Generating invocation graph for token of length %d bytes", len(tokenBytes))
	
	result, err := h.graph.GenerateInvocationGraph(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Invocation graph generation failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate invocation graph", err)
		return
	}

	log.Printf("[INFO] Successfully generated invocation graph: %d nodes, %d edges, is_invocation=%v", 
		len(result.Nodes), len(result.Edges), result.IsInvocation)
	respondJSON(w, http.StatusOK, result)
}

// GenerateInvocationGraphFile handles POST /api/graph/invocation/file
func (h *GraphHandler) GenerateInvocationGraphFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Invocation graph file generation request from %s", r.RemoteAddr)

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

	log.Printf("[DEBUG] Generating invocation graph for file %s (%d bytes)", 
		header.Filename, len(tokenBytes))
	
	result, err := h.graph.GenerateInvocationGraph(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Invocation graph generation failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate invocation graph", err)
		return
	}

	log.Printf("[INFO] Successfully generated invocation graph from file: %d nodes, %d edges, is_invocation=%v", 
		len(result.Nodes), len(result.Edges), result.IsInvocation)
	respondJSON(w, http.StatusOK, result)
}