package handlers

import (
	"encoding/json"
	"io"
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

// GenerateGraph handles JSON requests with token string
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

	// Auto-detect format and normalize to bytes
	tokenBytes, err := utils.NormalizeToken(req.Token, req.Format)
	if err != nil {
		log.Printf("[ERROR] Failed to normalize token: %v", err)
		respondError(w, http.StatusBadRequest, 
			"Invalid token format. Supported formats: base64, hex. "+
			"Or upload a .car file.", err)
		return
	}

	log.Printf("[DEBUG] Generating graph for token of length %d bytes", len(tokenBytes))
	result, err := h.graph.GenerateDelegationGraph(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Graph generation failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate graph", err)
		return
	}

	log.Printf("[INFO] Successfully generated graph: %d nodes, %d edges", 
		len(result.Nodes), len(result.Edges))
	respondJSON(w, http.StatusOK, result)
}

// GenerateGraphFile handles multipart file uploads
func (h *GraphHandler) GenerateGraphFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[INFO] Graph generation file upload request from %s", r.RemoteAddr)

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

	// Generate graph
	result, err := h.graph.GenerateDelegationGraph(tokenBytes)
	if err != nil {
		log.Printf("[ERROR] Graph generation failed: %v", err)
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate graph", err)
		return
	}

	log.Printf("[INFO] Successfully generated graph from file: %d nodes, %d edges", 
		len(result.Nodes), len(result.Edges))
	respondJSON(w, http.StatusOK, result)
}