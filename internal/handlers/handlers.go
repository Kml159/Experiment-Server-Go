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
	"experiment-server/internal/utils"
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
	client.ComputerAddress = utils.ReadUserIP(r)
	if err != nil {
		message := "Invalid request body"
		log.Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	contains := clients.Contains(client.ComputerAddress)
	if contains {
		log.Printf("This client already registered")
	} else {
		log.Printf("Appending client: %s", client.String())
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
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := clients.Get(utils.ReadUserIP(r))
	if err != nil {
		log.Printf("Unauthorized client: %s.", utils.ReadUserIP(r))
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	address := utils.ReadUserIP(r)
	clients.Activate(address)
	experiment := experiments.Subscribe()

	if experiment == nil {
		http.Error(w, "No experiments available", http.StatusNotFound)
		return
	}

	clients.AppendActiveExperiment(address, experiment.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*experiment)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	address := utils.ReadUserIP(r)

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client.ComputerAddress = address
	err = clients.Update(address, &client)
	if err != nil {
		log.Printf("Unauthorized client: %s.", address)
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	log.Printf("Status Updated: [%s]", address)
	w.WriteHeader(http.StatusOK)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var file output.Output
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received file: %s", file.FileName)

	parts := strings.Split(file.FileName, "_")
	if len(parts) < 2 {
		http.Error(w, "Invalid file name format", http.StatusBadRequest)
		return
	}

	ExperimentId := strings.TrimSuffix(strings.Split(parts[1], ".")[0], "")

	clientObj, err := clients.Get(utils.ReadUserIP(r))
	if err != nil {
		log.Printf("Unauthorized client: %s.", utils.ReadUserIP(r))
		http.Error(w, "Unauthorized client", http.StatusUnauthorized)
		return
	}

	for _, id := range clientObj.CompletedExperimentIDs {
		if id == ExperimentId {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "File already received"})
			return
		}
	}

	data, err := base64.StdEncoding.DecodeString(file.FileData)
	if err != nil {
		http.Error(w, "Failed to decode file data", http.StatusBadRequest)
		return
	}

	if err := os.MkdirAll(cfg.ReceivedOutputFilePath, 0755); err != nil {
		http.Error(w, "Failed to create output directory", http.StatusInternalServerError)
		return
	}

	outputPath := filepath.Join(cfg.ReceivedOutputFilePath, file.FileName)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	if err := clients.RemoveActiveExperiment(utils.ReadUserIP(r), ExperimentId); err != nil {
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
