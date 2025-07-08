package main

import (
	"experiment-server/internal/config"
	"experiment-server/internal/records/experiments"
	"experiment-server/internal/routes"
	"experiment-server/internal/workers/checker"
	"experiment-server/internal/workers/reporter"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	
	cfg := config.Load()

	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
	defer logFile.Close()

	experiments.Initialize(cfg)

	multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	router := routes.NewRouter()

	checker.Check()
	reporter.Report()
	
	log.Printf("Server starting on http://%s", "0.0.0.0:3754")
	log.Fatal(http.ListenAndServe("0.0.0.0:3754", router))
}
