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
	experiments  map[string]parameter.Parameter
	unsubscribed map[string]*parameter.Parameter
	completed 	 map[string]*parameter.Parameter
)

func init(){
	experiments = parameter.GenerateParamCombinations(duplicate)
	fmt.Println("Generated experiment parameters:")
	for key, p := range experiments{
		unsubscribed[key] = &p
		p.Print()
	}
	unsubscribed = make(map[string]*parameter.Parameter, len(experiments))
	completed = make(map[string]*parameter.Parameter, len(experiments))
}

func Add(param parameter.Parameter) error {
	mu.Lock()
	defer mu.Unlock()
	_, ok := experiments[param.ID]
	if !ok{
		experiments[param.ID] = param
	}
	unsubscribed[param.ID] = &param
	return nil
}

func Subcribe() *parameter.Parameter { 
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