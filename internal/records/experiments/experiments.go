package experiments

import (
	"experiment-server/internal/config"
	"experiment-server/internal/models/parameter"
	"fmt"
	"log"
	"sync"
)


var (
	mu           sync.Mutex
	experiments  map[string]parameter.Parameter
	unsubscribed map[string]*parameter.Parameter
	completed    map[string]*parameter.Parameter
)

func Initialize(cfg *config.Config) {
	experiments = parameter.GenerateParamCombinations(cfg.ExperimentDuplicate, cfg)

	if cfg.SubtractCompletedExperiments {
		err := parameter.SubtractCompleted(&experiments, cfg.ReceivedOutputFilePath, cfg.ExperimentBaseId)
		if err != nil {
			log.Printf("Error subtracting completed experiments: %v", err)
		}
	}

	for experiment := range experiments {
		log.Println(experiment)
	}

	unsubscribed = make(map[string]*parameter.Parameter, len(experiments))
	completed = make(map[string]*parameter.Parameter, len(experiments))
	log.Printf("Generated experiment parameters:")
	for key, p := range experiments {
		unsubscribed[key] = &p
		p.Print()
	}
	fmt.Print("\n")
}

func UnsubscribedCount() int {
	mu.Lock()
	defer mu.Unlock()
	return len(unsubscribed)
}

func Add(params ...parameter.Parameter) {
	mu.Lock()
	defer mu.Unlock()
	for _, param := range params {
		if _, ok := experiments[param.ID]; !ok {
			experiments[param.ID] = param
		}
		unsubscribed[param.ID] = &param
	}
}

func Subscribe() *parameter.Parameter {
	mu.Lock()
	defer mu.Unlock()
	for id, param := range unsubscribed {
		delete(unsubscribed, id)
		return param
	}
	return nil
}

func Completed(key string) error {
	mu.Lock()
	defer mu.Unlock()

	exp, ok := experiments[key]
	if !ok {
		return fmt.Errorf("experiment %q not found", key)
	}
	completed[key] = &exp
	return nil
}

func IsDone() bool {
	mu.Lock()
	defer mu.Unlock()
	return len(completed) >= len(experiments)
}

func Parameters(ids ...string) ([]parameter.Parameter, error) {
	mu.Lock()
	defer mu.Unlock()

	params := make([]parameter.Parameter, 0, len(ids))
	for _, id := range ids {
		exp, ok := experiments[id]
		if !ok {
			return nil, fmt.Errorf("experiment %q not found", id)
		}
		params = append(params, exp)
	}
	return params, nil
}
