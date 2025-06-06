package experiments

import (
	"experiment-server/internal/models/parameter"
	"fmt"
	"sync"
)

const (
	duplicate = 5
)

var (
	mu           sync.Mutex
	unsubscribed = make(map[string]parameter.Parameter)
	experiments  = make(map[string]parameter.Parameter)
)

func Add(param parameter.Parameter) error {
	mu.Lock()
	defer mu.Unlock()
	unsubscribed[param.ID] = param
	return nil
}

func init(){
	experiments = parameter.GenerateParamCombinations(duplicate)
	fmt.Println("Generated experiment parameters:")
	for _, p := range experiments{
		p.Print()
	}
	unsubscribed = experiments
}

func Subcribe() *parameter.Parameter {
    mu.Lock()
    defer mu.Unlock()
    for id, param := range unsubscribed {
        delete(unsubscribed, id)
        return &param
    }
    return nil 
}