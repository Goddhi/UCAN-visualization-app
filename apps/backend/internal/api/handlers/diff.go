package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/goddhi/ucan-visualizer/internal/models"
	"github.com/goddhi/ucan-visualizer/internal/services/diff"
)

type DiffHandler struct {
	service *diff.Service
}

func NewDiffHandler(svc *diff.Service) *DiffHandler {
	return &DiffHandler{
		service: svc,
	}
}

func (h *DiffHandler) GenerateDiff(w http.ResponseWriter, r *http.Request) {
	var req models.DiffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	diffs, err := h.service.GenerateDiff(req.ParentToken, req.ChildToken)
	if err != nil {
		respondError(w, http.StatusUnprocessableEntity, "Failed to generate diff", err)
		return
	}

	respondJSON(w, http.StatusOK, models.DiffResponse{Diffs: diffs})
}