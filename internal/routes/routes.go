package routes

import (
	"log"
	"net/http"

	"experiment-server/internal/config"
	"experiment-server/internal/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQ] %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func NewRouter() http.Handler {

	cfg := config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, cfg)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, cfg)
	})
	mux.HandleFunc("/get_experiment", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetExperimentHandler(w, r, cfg)
	})
	mux.HandleFunc("/update_status", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateStatusHandler(w, r, cfg)
	})
	mux.HandleFunc("/upload_file", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadFileHandler(w, r, cfg)
	})

	return loggingMiddleware(mux)
}
