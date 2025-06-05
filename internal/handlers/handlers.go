package handlers

import (
	"encoding/json"
	"experiment-server/internal/config"
	"experiment-server/internal/models/client"
	"experiment-server/internal/records/clients"
	"experiment-server/internal/dto/register"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World! This is the Go net/http server template.")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	config := config.Load()

	if r.Method != http.MethodPost {
		message := "Method Not Allowed"
		http.Error(w, message, http.StatusMethodNotAllowed)
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Appending client: ", client.String())
	clients.Add(&client)

	response := register.RegisterResponse{
		Status:       "Successful",
		ThreadAmount: config.ThreadAmount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetExperimentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	client, err := clients.Get(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Client is not registered", http.StatusExpectationFailed)
	}

	clients.Activate(r.RemoteAddr)

	// TODO
}
