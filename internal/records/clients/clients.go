package clients

import (
	"experiment-server/internal/models/client"
	"fmt"
	"sync"
)

var (
	mu        sync.Mutex
	byAddress = make(map[string]*client.Client)
	total     int
	active    int
)

func Add(c *client.Client) error {
	mu.Lock()
	defer mu.Unlock()
	byAddress[c.ComputerAddress] = c
	total++
	if c.IsActive {
		active++
	}
	return nil
}

func Get(key string) (*client.Client, error) {
	mu.Lock()
	defer mu.Unlock()
	c, ok := byAddress[key]
	if !ok {
		return nil, fmt.Errorf("client not found: %s", key)
	}
	return c, nil
}

func Stats() {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("Total Clients:", total, "\tActive:", active, "\tInactive:", total-active)
}

// TODO: fix
func Deactivate(key string) error {
	mu.Lock()
	defer mu.Unlock()
	c, ok := byAddress[key]
	if !ok {
		return fmt.Errorf("client not found: %s", key)
	}
	if !c.IsActive {
		return fmt.Errorf("client is already deactive: %s", key)
	}
	c.IsActive = false
	active--
	return nil
}

func Activate(key string) error {
	mu.Lock()
	defer mu.Unlock()
	c, ok := byAddress[key]
	if !ok {
		return fmt.Errorf("client not found: %s", key)
	}
	if c.IsActive {
		return nil
	}
	c.IsActive = true
	active++
	return nil
}

func Contains(key string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, ok := byAddress[key]
	return ok
}

func AppendActiveExperiment(ClientAddress string, ExperimentId string) error {
	mu.Lock()
	defer mu.Unlock()
	client, ok := byAddress[ClientAddress]
	if !ok{
		return fmt.Errorf("Client Address does not exists: " + ClientAddress)
	}
	client.ActiveExperimentIDs = append(client.ActiveExperimentIDs, ExperimentId)
	return nil
}

func RemoveActiveExperiment(ClientAddress string, ExperimentId string) error {
    mu.Lock()
    defer mu.Unlock()
    client, ok := byAddress[ClientAddress]
    if !ok {
        return fmt.Errorf("client Address does not exist: %s", ClientAddress)
    }
    ids := client.ActiveExperimentIDs
    for i, id := range ids {
        if id == ExperimentId {
            client.ActiveExperimentIDs = append(ids[:i], ids[i+1:]...)
			client.CompletedExperimentIDs = append(client.CompletedExperimentIDs, ExperimentId)
            return nil
        }
    }
    return fmt.Errorf("experiment ID not found: %s", ExperimentId)
}