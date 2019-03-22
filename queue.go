package queue

import (
	"sync"
)

type Queues struct {
	limit    int
	next     int
	choosing sync.Mutex
	list     []*sync.Mutex
}

type queuesControl interface {
	Init()
	Get()  sync.Mutex
}

func New(limit int) *Queues {
	que := new(Queues)
	que.Init(limit)
	return que
}

func (q *Queues) Init(limit int) {
	if limit <= 0 {
		panic("Queues: invalid limit")
	}

	q.next  = limit
	q.limit = limit
	q.list = make([]*sync.Mutex, limit)

	for i := 0; i < limit; i++ {
		q.list[i] = new(sync.Mutex)
	}
}

func (q *Queues) Get() *sync.Mutex {
	q.choosing.Lock()
	door := getNextDoor(q)
	q.choosing.Unlock()
	return door
}

func getNextDoor(q *Queues) *sync.Mutex {
	q.next++
	if q.next >= q.limit {
		q.next = 0
	}
	return q.list[q.next]
}
