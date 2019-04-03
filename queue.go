package queue

import (
	"sync"
)

func New(limit int) queuesControl {
	que := new(queues)
	initQueues(que, limit)
	return que
}

func initQueues(q *queues, limit int) {
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

func (q *queues) Get() *sync.Mutex {
	q.choosing.Lock()
	queue := getNextQueue(q)
	q.choosing.Unlock()
	return queue
}

func getNextQueue(q *queues) *sync.Mutex {
	q.next++
	if q.next >= q.limit {
		q.next = 0
	}
	return q.list[q.next]
}

func NewMult(multLimit Mult) multQueuesControl {
	m := new(multQueues)
	m.mult = make(queueMap)

	for key, limit := range multLimit {
		m.mult[key] = New(limit)
	}

	return m
}

func (m *multQueues) GetGroup(group string) *sync.Mutex {
	m.choosing.Lock()
	q, exist := m.mult[group]
	m.choosing.Unlock()

	if !exist {
		panic("Queues: invalid queue group")
	}
	return q.Get()
}
