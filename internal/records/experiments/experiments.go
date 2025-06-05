package experiments

import (
	"experiment-server/internal/models/parameter"
	"sync"
)

var (
	mu           sync.Mutex
	unsubscribed = make(map[string]parameter.Parameter)
)

func Add(param parameter.Parameter) error {
	mu.Lock()
	defer mu.Unlock()
	unsubscribed[param.ID] = param
	return nil
}
