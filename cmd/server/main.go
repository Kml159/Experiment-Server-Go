package main

import (
	"log"
	"net/http"

	"experiment-server/internal/routes"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
    log.Println("Verbose logging ENABLED")

    router := routes.NewRouter()
    log.Printf("Server starting on http://%s", "localhost:8080")
    log.Fatal(http.ListenAndServe("localhost:8080", router))
}