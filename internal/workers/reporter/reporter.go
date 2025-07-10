package reporter

import (
	"log"
	"time"

	"experiment-server/internal/records/clients"
	"experiment-server/internal/records/experiments"
)

const (
	tick = 15
)

func Report() {
	go func() {
		ticker := time.NewTicker(tick * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			all := clients.Count()
			active := clients.ActiveClientCount()
			lost := all - active

			unsubscribedExperiments := experiments.UnsubscribedCount()
			allExperiments := experiments.AllCount()
			log.Printf("Clients: [%dT: %dA, %dL], Experiments: [%dT: %dU, %dC]\n",
				all,
				active,
				lost,
				allExperiments,
				unsubscribedExperiments,
				experiments.CompletedCount(),
			)
		}
	}()
}
