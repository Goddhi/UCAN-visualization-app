package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/goddhi/ucan-visualizer/internal/api"
	"github.com/goddhi/ucan-visualizer/internal/config"
)

func main() {
	cfg := config.Load()

	handler := api.SetupRouter()

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}