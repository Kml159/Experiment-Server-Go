package clients

import (
	"experiment-server/internal/models/client"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	mu        sync.Mutex
	clientMap = make(map[string]*client.Client)
	total     int
	active    int
)

func Add(c *client.Client) error {
	mu.Lock()
	defer mu.Unlock()
	clientMap[c.ComputerName] = c
	total++
	if c.IsActive {
		active++
	}
	return nil
}

func Get(key string) (*client.Client, error) {
	mu.Lock()
	defer mu.Unlock()
	c, ok := clientMap[key]
	if !ok {
		return nil, fmt.Errorf("client not found: %s", key)
	}
	return c, nil
}

func Stats() {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Total Clients: %d\tActive: %d\tInactive: %d", total, active, total-active)
}

func Count() int {
	mu.Lock()
	defer mu.Unlock()
	return len(clientMap)
}

func Deactivate(key string) error {
	mu.Lock()
	defer mu.Unlock()
	c, ok := clientMap[key]
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
	c, ok := clientMap[key]
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
	_, ok := clientMap[key]
	return ok
}

func AppendActiveExperiment(key string, ExperimentId string) error {
	mu.Lock()
	defer mu.Unlock()
	client, ok := clientMap[key]
	if !ok {
		return fmt.Errorf("Client Address does not exists: " + key)
	}
	client.ActiveExperimentIDs = append(client.ActiveExperimentIDs, ExperimentId)
	return nil
}

func RemoveActiveExperiment(key string, ExperimentId string) error {
	mu.Lock()
	defer mu.Unlock()
	client, ok := clientMap[key]
	if !ok {
		return fmt.Errorf("client Address does not exist: %s", key)
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

func Update(key string, client *client.Client) error {
	mu.Lock()
	defer mu.Unlock()
	_, ok := clientMap[key]
	if !ok {
		return fmt.Errorf("client Address does not exist: %s", key)
	}
	client.LastStatusReceived = time.Now()
	clientMap[key] = client
	return nil
}

func Clients() []client.Client {
	mu.Lock()
	defer mu.Unlock()
	clients := make([]client.Client, 0, len(clientMap))
	for _, c := range clientMap {
		clients = append(clients, *c)
	}
	return clients
}

func ActiveClientCount() int {
	mu.Lock()
	defer mu.Unlock()
	return active
}
