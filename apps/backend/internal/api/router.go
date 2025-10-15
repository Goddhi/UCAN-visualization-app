package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/goddhi/ucan-visualizer/internal/api/handlers"
)

func SetupRouter() http.Handler {
	r := mux.NewRouter()

	// Initialize handlers
	parseHandler := handlers.NewParseHandler()
	validateHandler := handlers.NewValidateHandler()
	graphHandler := handlers.NewGraphHandler()

	// Health check
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Parse endpoints
	api.HandleFunc("/parse/delegation", parseHandler.ParseDelegation).Methods("POST")
	api.HandleFunc("/parse/delegation/file", parseHandler.ParseFile).Methods("POST")
	
	// Validate endpoints
	api.HandleFunc("/validate/chain", validateHandler.ValidateChain).Methods("POST")
	api.HandleFunc("/validate/chain/file", validateHandler.ValidateFile).Methods("POST")
	
	// Graph endpoints
	api.HandleFunc("/graph/delegation", graphHandler.GenerateGraph).Methods("POST")
	api.HandleFunc("/graph/delegation/file", graphHandler.GenerateGraphFile).Methods("POST")

	return r
}