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

type queueMap map[string]queuesControl

type multQueues struct {
	choosing sync.Mutex
	mult     queueMap
}

type multQueuesControl interface {
	GetGroup(string) *sync.Mutex
}
