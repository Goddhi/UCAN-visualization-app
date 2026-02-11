package api

import (
	"net/http"

	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/goddhi/ucan-visualizer/internal/api/handlers"
	"github.com/goddhi/ucan-visualizer/internal/services/diff" 
)

func SetupRouter() http.Handler {
	r := mux.NewRouter()

	diffService := diff.NewService() 

	// Initialize handlers
	parseHandler := handlers.NewParseHandler()
	validateHandler := handlers.NewValidateHandler()
	graphHandler := handlers.NewGraphHandler()
	diffHandler := handlers.NewDiffHandler(diffService) 

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

	// diff endpoint 
	api.HandleFunc("/diff", diffHandler.GenerateDiff).Methods("POST")
	
	cors := gorillahandlers.CORS(
		gorillahandlers.AllowedOrigins([]string{"*"}),
		gorillahandlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		gorillahandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	return cors(r)
}
