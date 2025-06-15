package checker

import (
	"experiment-server/internal/records/clients"
	"experiment-server/internal/records/experiments"
	"fmt"
	"time"
)

const (
	checkerInterval            = 60
	tolerance                  = 2
	clientStatusSenderInterval = 60 * 10
)

func Check() {
	go func() {
		ticker := time.NewTicker(checkerInterval * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			for _, client := range clients.Clients() {
				if client.LastStatusReceived.IsZero() {
					continue
				}

				if time.Since(client.LastStatusReceived) > time.Duration(clientStatusSenderInterval*tolerance)*time.Second {
					fmt.Println("\nClient", client.ComputerName, "is unresponsive. Deactivating from clients list.")

					err := clients.Deactivate(client.ComputerName)
					if err != nil {
						fmt.Printf("error deactivating client %s: %v\n", client.ComputerName, err)
						continue
					}

					lostParams, err := experiments.Parameters(client.ActiveExperimentIDs...)
					if err != nil {
						fmt.Printf("error getting experiment parameters for client %s: %v\n", client.ComputerName, err)
						continue
					}

					experiments.Add(lostParams...)
				}
			}
		}
	}()
}
