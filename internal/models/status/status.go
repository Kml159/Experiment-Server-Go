package status

import (
	"fmt"
	"time"
)

type Status struct{
	StartTime time.Time
}

func (s Status) GetUpTime() time.Duration {
    return time.Since(s.StartTime)
}

func (s Status) Print() {
    fmt.Printf("Status{StartTime: %s, UpTime: %s}\n", s.StartTime.Format("2006-01-02 15:04:05"), s.GetUpTime())
}

var ServerStatus Status

func init(){
	ServerStatus = Status{StartTime: time.Now()}
}