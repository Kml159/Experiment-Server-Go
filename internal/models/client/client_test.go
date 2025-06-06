package client

import (
    "testing"
    "time"
)

func TestNewClient(t *testing.T) {
    c := NewClient()
    if c.IsExperimentRunning {
        t.Errorf("expected IsExperimentRunning to be false, got true")
    }
    if c.IsActive {
        t.Errorf("expected IsActive to be false, got true")
    }
    if len(c.CompletedExperimentIDs) != 0 {
        t.Errorf("expected CompletedExperimentIDs to be empty, got %v", c.CompletedExperimentIDs)
    }
    if c.ComputerName != "" {
        t.Errorf("expected ComputerName to be empty, got %s", c.ComputerName)
    }
    if c.ComputerAddress != "" {
        t.Errorf("expected ComputerAddress to be empty, got %s", c.ComputerAddress)
    }
    if c.StartTime.IsZero() {
        t.Errorf("expected StartTime to be set, got zero value")
    }
    if c.LastStatusReceived.IsZero() {
        t.Errorf("expected LastStatusReceived to be set, got zero value")
    }
    if c.ExperimentStartTime.IsZero() {
        t.Errorf("expected ExperimentStartTime to be set, got zero value")
    }
}

func TestClientString(t *testing.T) {
    c := NewClient()
    str := c.String()
    if str == "" {
        t.Errorf("expected non-empty string from String(), got empty string")
    }
    if want := "Client("; str[:7] != want {
        t.Errorf("expected string to start with %q, got %q", want, str[:7])
    }
}

func TestClientTimeFields(t *testing.T) {
    c := NewClient()
    now := time.Now()
    if c.StartTime.After(now) {
        t.Errorf("expected StartTime to be before or equal to now, got %v", c.StartTime)
    }
}