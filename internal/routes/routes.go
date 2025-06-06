package routes

import (
	"net/http"

	"experiment-server/internal/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/get_experiment", handlers.GetExperimentHandler)
	mux.HandleFunc("/update_status", handlers.UpdateStatusHandler)
	mux.HandleFunc("/upload_file", handlers.UploadFileHandler)
	return mux
}
