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
		fmt.Println(message)
		http.Error(w, message, http.StatusMethodNotAllowed)
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		message := "Invalid request body"
		fmt.Println(message)
		http.Error(w, message, http.StatusBadRequest)
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

	_, err := clients.Get(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Client is not registered", http.StatusExpectationFailed)
	}

	clients.Activate(r.RemoteAddr)
	experiment := experiments.Subcribe()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*experiment)
}

func UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	var client client.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client.ComputerAddress = r.RemoteAddr
	clients.Update(r.RemoteAddr, &client)
	
    fmt.Println("[" + time.Now().Format(time.RFC3339) + "]", "Status Updated: [" + r.RemoteAddr + "]", )
}

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

    fmt.Println("Received file:", file.FileName)

    parts := strings.Split(file.FileName, "_")
    if len(parts) < 2 {
        http.Error(w, "Invalid file name format", http.StatusBadRequest)
        return
    }

    ExperimentId := strings.TrimSuffix(strings.Split(parts[1], ".")[0], "")

    clientObj, err := clients.Get(r.RemoteAddr)
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

    if err := clients.RemoveActiveExperiment(r.RemoteAddr, ExperimentId); err != nil {
        fmt.Printf("Error moving experiment ID %s: %v\n", ExperimentId, err)
    }

	if err := experiments.Completed(ExperimentId); err != nil{
		fmt.Println("Error on marking experiment as done!")
	}

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "File received"})
}