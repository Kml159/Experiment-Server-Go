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
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	config := config.Load()

	if r.Method != http.MethodPost {
		message := "Method Not Allowed"
		log.Println(message)
		http.Error(w, message, http.StatusMethodNotAllowed)
		return
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	client.ComputerAddress = strings.Split(r.RemoteAddr, ":")[0]
	if err != nil {
		message := "Invalid request body"
		log.Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	log.Printf("Appending client: %s", client.String())
	clients.Add(&client)

	response := register.RegisterResponse{
		Status:                          "Successful",
		ThreadAmount:                    config.ThreadAmount,
		ClientSendUpdateStatusInSeconds: config.ClientSendUpdateStatusInSeconds,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetExperimentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := clients.Get(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		http.Error(w, "Client is not registered", http.StatusUnauthorized)
		log.Printf("Unauthorized access attempt from %s: %v\n", r.RemoteAddr, err)
		return
	}

	address := strings.Split(r.RemoteAddr, ":")[0]
	clients.Activate(address)
	experiment := experiments.Subcribe()

	if experiment == nil {
		http.Error(w, "No experiments available", http.StatusNotFound)
		return
	}

	clients.AppendActiveExperiment(address, experiment.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*experiment)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client.ComputerAddress = strings.Split(r.RemoteAddr, ":")[0]
	clients.Update(strings.Split(r.RemoteAddr, ":")[0], &client)

	log.Printf("[%s] Status Updated: [%s]", time.Now().Format(time.RFC3339), r.RemoteAddr)}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
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

	clientObj, err := clients.Get(strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		http.Error(w, "Client not registered", http.StatusUnauthorized)
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

	if err := os.MkdirAll("received_output", 0755); err != nil {
		http.Error(w, "Failed to create output directory", http.StatusInternalServerError)
		return
	}

	outputPath := filepath.Join("received_output", file.FileName)
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	if err := clients.RemoveActiveExperiment(strings.Split(r.RemoteAddr, ":")[0], ExperimentId); err != nil {
		log.Printf("Error moving experiment ID %s: %v\n", ExperimentId, err)
	}

	if err := experiments.Completed(ExperimentId); err != nil {
		log.Printf("Error on marking experiment as done!")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "File received"})
}
