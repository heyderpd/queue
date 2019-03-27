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
