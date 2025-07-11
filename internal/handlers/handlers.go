package handlers

import (
	"encoding/base64"
	"encoding/json"
	"experiment-server/internal/config"
	"experiment-server/internal/dto/register"
	"experiment-server/internal/models/client"
	"experiment-server/internal/models/output"
	"experiment-server/internal/records/clients"
	"experiment-server/internal/records/experiments"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	fmt.Fprintln(w, "Hello, World!")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {

	if r.Method != http.MethodPost {
		message := "Method Not Allowed"
		log.Println(message)
		http.Error(w, message, http.StatusMethodNotAllowed)
		return
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	client.ComputerAddress = r.RemoteAddr
	if err != nil {
		message := "Invalid request body"
		log.Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	if clients.Contains(client.ComputerName) {
		log.Printf("This client already registered")
	} else {
		log.Printf("Appending new client: %s", client.String())
		clients.Add(&client)
	}

	response := register.RegisterResponse{
		Status:                          "Successful",
		ThreadAmount:                    cfg.ThreadAmount,
		ClientSendUpdateStatusInSeconds: cfg.ClientSendUpdateStatusInSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetExperimentHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {

	if r.Method != http.MethodGet {
		log.Printf("Method Not Allowed: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := clients.Get(r.Header.Get("ComputerName"))
	if err != nil {
		log.Printf("Unauthorized client: %s.", r.RemoteAddr)
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	ComputerName := r.Header.Get("ComputerName")
	clients.Activate(ComputerName)
	experiment := experiments.Subscribe()

	if experiment == nil {
		log.Printf("No experiments available for client: %s", r.Header.Get("ComputerName"))
		http.Error(w, "No experiments available", http.StatusNotFound)
		return
	}

	clients.AppendActiveExperiment(ComputerName, experiment.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*experiment)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {

	if r.Method != http.MethodPost {
		log.Printf("Method Not Allowed: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client.ComputerAddress = r.RemoteAddr
	err = clients.Update(client.ComputerName, &client)
	if err != nil {
		log.Printf("Unauthorized client: %s.", client.ComputerName)
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	log.Printf("Status Updated: [%s]", client.ComputerName)
	w.WriteHeader(http.StatusOK)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	if r.Method != http.MethodPost {
		log.Printf("Method Not Allowed: %s", r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var file output.Output
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received file: %s", file.FileName)

	parts := strings.Split(file.FileName, "_")
	if len(parts) < 2 {
		log.Printf("Invalid file name format: %s", file.FileName)
		http.Error(w, "Invalid file name format", http.StatusBadRequest)
		return
	}

	ExperimentId := strings.TrimSuffix(strings.Split(parts[1], ".")[0], "")

	clientObj, err := clients.Get(r.Header.Get("ComputerName"))
	if err != nil {
		log.Printf("Unauthorized client: %s.", r.RemoteAddr)
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	for _, id := range clientObj.CompletedExperimentIDs {
		if id == ExperimentId {
			w.WriteHeader(http.StatusOK)
			log.Printf("File already received for experiment ID: %s from client: %s", ExperimentId, r.Header.Get("ComputerName"))
			json.NewEncoder(w).Encode(map[string]string{"status": "File already received"})
			return
		}
	}

	data, err := base64.StdEncoding.DecodeString(file.FileData)
	if err != nil {
		log.Printf("Failed to decode file data for file %s: %v", file.FileName, err)
		http.Error(w, "Failed to decode file data", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(cfg.ReceivedOutputFilePath, 0755); err != nil {
		log.Printf("Failed to create output directory %s: %v", cfg.ReceivedOutputFilePath, err)
		http.Error(w, "Failed to create output directory", http.StatusInternalServerError)
		return
	}

	outputPath := filepath.Join(cfg.ReceivedOutputFilePath, file.FileName)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Printf("Failed to write file %s: %v", outputPath, err)
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	if err := clients.RemoveActiveExperiment(r.Header.Get("ComputerName"), ExperimentId); err != nil {
		log.Printf("Error moving experiment ID %s: %v\n", ExperimentId, err)
		http.Error(w, "Error moving experiment ID: Requesting reregister", http.StatusUnauthorized)
		return
	}

	if err := experiments.Completed(ExperimentId); err != nil {
		log.Printf("Error on marking experiment as done!")
		http.Error(w, "Error on marking experiment as done!: Requesting reregister", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "File received"})
}
