package main

import (
	"log"
	"net/http"
	"os"

	"experiment-server/internal/routes"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Println("Verbose logging ENABLED")

	router := routes.NewRouter()
	log.Printf("Server starting on http://%s", "localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
