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
			log.Printf("[%s] Clients: %d, Unsubscribed Experiments: %d\n",
				time.Now().Format("2006-01-02 15:04:05"),
				clients.Count(),
				experiments.UnsubscribedCount(),
			)
		}
	}()
}
