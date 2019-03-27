package queue

import (
	"sync"
)

type queues struct {
	limit    int
	next     int
	choosing sync.Mutex
	list     []*sync.Mutex
}

type queuesControl interface {
	Get() *sync.Mutex
}

type Mult map[string]int

type multQueues map[string]queuesControl

type multQueuesControl interface {
	Get(string) *sync.Mutex
}
