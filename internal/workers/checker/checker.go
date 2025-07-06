package checker

import (
	"experiment-server/internal/config"
	"experiment-server/internal/records/clients"
	"experiment-server/internal/records/experiments"
	"fmt"
	"time"
)

const (
	checkerInterval            = 30
	tolerance                  = 2
)

func Check() {

	config := config.Load()
	go func() {
		ticker := time.NewTicker(checkerInterval * time.Second)
		defer ticker.Stop()

		for range ticker.C {

			for _, client := range clients.Clients() {

				if client.LastStatusReceived.IsZero() {
					continue
				}

				if !client.IsActive {
					continue
				}

				if time.Since(client.LastStatusReceived) > time.Duration(config.ClientSendUpdateStatusInSeconds*tolerance)*time.Second {
					fmt.Println("\nClient", client.ComputerName, "is unresponsive. Deactivating from clients list.")

					err := clients.Deactivate(client.ComputerAddress)
					if err != nil {
						fmt.Printf("error deactivating client [%s]: %v\n", client.ComputerAddress, err)
						continue
					}

					lostParams, err := experiments.Parameters(client.ActiveExperimentIDs...)
					if err != nil {
						fmt.Printf("error getting experiment parameters for client %s: %v\n", client.ComputerName, err)
						continue
					}
					
					fmt.Printf("Lost parameters added: %+v\n", lostParams)
					experiments.Add(lostParams...)
				}
			}
		}
	}()
}
