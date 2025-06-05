package main

import (
	"log"
	"net/http"

	"experiment-server/internal/config"
	"experiment-server/internal/models/status"
	"experiment-server/internal/routes"
)

func main() {
	cfg := config.Load()
	router := routes.NewRouter()
	status.ServerStatus.Print()

	log.Printf("Server starting on %s...", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
