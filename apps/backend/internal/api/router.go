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

	r.HandleFunc("/", handlers.RootHandler).Methods("GET")

	// Health check
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/healthz", handlers.HealthCheck).Methods("GET")

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	
	// Parse endpoints
	api.HandleFunc("/parse/delegation", parseHandler.ParseDelegation).Methods("POST")
	api.HandleFunc("/parse/delegation/file", parseHandler.ParseFile).Methods("POST")
	api.HandleFunc("/parse/chain", parseHandler.ParseChain).Methods("POST")
	api.HandleFunc("/parse/chain/file", parseHandler.ParseChainFile).Methods("POST")
	api.HandleFunc("/parse/invocation", parseHandler.ParseInvocation).Methods("POST")
	api.HandleFunc("/parse/invocation/file", parseHandler.ParseInvocationFile).Methods("POST")
	
	// Validate endpoints
	api.HandleFunc("/validate/chain", validateHandler.ValidateChain).Methods("POST")
	api.HandleFunc("/validate/chain/file", validateHandler.ValidateFile).Methods("POST")
	
	// Graph endpoints  
	api.HandleFunc("/graph/delegation", graphHandler.GenerateGraph).Methods("POST")
	api.HandleFunc("/graph/delegation/file", graphHandler.GenerateGraphFile).Methods("POST")
	api.HandleFunc("/graph/invocation", graphHandler.GenerateInvocationGraph).Methods("POST")
	api.HandleFunc("/graph/invocation/file", graphHandler.GenerateInvocationGraphFile).Methods("POST")

	return r
}