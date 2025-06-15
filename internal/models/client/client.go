package client

import (
	"fmt"
	"time"
)

type Client struct {
	IsExperimentRunning    bool      `json:"is_experiment_running"`
	IsActive               bool      `json:"is_active"`
	CompletedExperimentIDs []string  `json:"completed_experiment_ids"`
	StartTime              time.Time `json:"start_time"`
	ActiveExperimentIDs    []string  `json:"active_experiment_ids"`
	ComputerName           string    `json:"computer_name"`
	ComputerAddress        string    `json:"computer_address"`
	LastStatusReceived     time.Time `json:"last_status_received"`
	ExperimentStartTime    time.Time `json:"experiment_start_time"`
	ID                     string    `json:"id"`
}

func (c Client) Print() {
	fmt.Println(c.String())
}

func NewClient() Client {
	now := time.Now()
	return Client{
		IsExperimentRunning:    false,
		IsActive:               false,
		CompletedExperimentIDs: []string{},
		StartTime:              now,
		ActiveExperimentIDs:    []string{},
		ComputerName:           "",
		ComputerAddress:        "",
		LastStatusReceived:     time.Time{},
		ExperimentStartTime:    now,
	}
}

func (c Client) String() string {
	return fmt.Sprintf(
		"Client(IsExperimentRunning=%v, IsActive=%v, CompletedExperimentIDs=%v, StartTime=%v, ActiveExperimentIDs=%v, ComputerName=%s, ComputerAddress=%s)",
		c.IsExperimentRunning, c.IsActive, c.CompletedExperimentIDs, c.StartTime.Format("2006-01-02 15:04:05"), c.ActiveExperimentIDs, c.ComputerName, c.ComputerAddress,
	)
}
