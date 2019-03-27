package queue

import (
	"sync"
	"fmt"
	"reflect"
)

func New(args interface{}) queuesControl {
	fmt.Println(args)
	fmt.Println(reflect.TypeOf(args))
	if val, ok := args.(int); ok {
		return newSimple(val)
	}
	if val, ok := args.(Mult); ok {
		return newMult(val)
	}
	panic("Queues: invalid argument")
}

func newMult(limit Mult) multQueuesControl {
	que := new(queues)
	initQueues(que, 10)
	return que
}

func newSimple(limit int) queuesControl {
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
	door := getNextDoor(q)
	q.choosing.Unlock()
	return door
}

func getNextDoor(q *queues) *sync.Mutex {
	q.next++
	if q.next >= q.limit {
		q.next = 0
	}
	return q.list[q.next]
}
