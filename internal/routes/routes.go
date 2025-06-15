package routes

import (
	"log"
	"net/http"

	"experiment-server/internal/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQ] %s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/get_experiment", handlers.GetExperimentHandler)
	mux.HandleFunc("/update_status", handlers.UpdateStatusHandler)
	mux.HandleFunc("/upload_file", handlers.UploadFileHandler)

	return loggingMiddleware(mux)
}
