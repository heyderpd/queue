package queue

import (
	"testing"
	"time"
	// "fmt"
)

var (
	sleepWeight = 5
	sleep = time.Millisecond * time.Duration(sleepWeight)
)

type jobTest func(chan bool)

func needPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}

func dontPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("Result in panic")
	}
}

func getNow() float64 {
	return float64(time.Now().UnixNano())
}

func getExpected(t int) float64 {
	return float64(t *1000 *1000 *sleepWeight)
}

func hasMuchTimeDiff(err float64, t *testing.T) {
	if 0.8 > err || err > 1.2 {
		t.Errorf("error: has so much time diff %v", err)
	}
}

func executeAssyncTestGeneric(t *testing.T, job jobTest, limit int, test int) {
	channel := make(chan bool)
	count := test

	base := float64(limit)
	if base > float64(test) {
		base = float64(test)
	}
	expected := getExpected(test) / base
	start := float64(0)

	go func(){
		time.Sleep(time.Millisecond * 1)
		start = getNow()

		for range make([]bool, test) {
			go job(channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDiff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}

func executeAssyncTest(t *testing.T, que queuesControl, limit int, test int) {
	job := func(channel chan bool) {
		// get mutex
		mutex := que.Get()
		mutex.Lock()
		defer mutex.Unlock()

		// do some work
		time.Sleep(sleep)
		channel <- true
	}
	executeAssyncTestGeneric(t, job, limit, test)
}

func executeAssyncTestWithMult(t *testing.T, que multQueuesControl, name string, limit int, test int) {
	job := func(channel chan bool) {
		// get mutex
		mutex := que.GetGroup(name)
		mutex.Lock()
		defer mutex.Unlock()

		// do some work
		time.Sleep(sleep)
		channel <- true
	}
	executeAssyncTestGeneric(t, job, limit, test)
}

func TestInitWithPanic(t *testing.T) {
	defer needPanic(t)

	que := New(0)
	que.Get()
}

func TestInit(t *testing.T) {
	defer dontPanic(t)

	que := New(1)
	que = New(5)
	que = New(10)
	que = New(200)
	que = New(3000)
	que.Get()
}

func TestBasic(t *testing.T) {
	defer dontPanic(t)

	que := New(1)
	mutex := que.Get()
	mutex.Lock()
	mutex.Unlock()
}

func TestSequencial(t *testing.T) {
	defer dontPanic(t)
	que := New(1)
	executeAssyncTest(t, que, 1, 10)
}

func TestSequencialWihtRoutines(t *testing.T) {
	defer dontPanic(t)
	que := New(1)
	executeAssyncTest(t, que, 1, 10)
}

func TestWithTwoQueue(t *testing.T) {
	defer dontPanic(t)
	que := New(2)
	executeAssyncTest(t, que, 2, 10)
}

func TestWithFiveQueue(t *testing.T) {
	defer dontPanic(t)
	que := New(5)
	executeAssyncTest(t, que, 5, 10)
}

func TestWithTeenQueue(t *testing.T) {
	defer dontPanic(t)
	que := New(10)
	executeAssyncTest(t, que, 10, 10)
}

func TestMoreQueueThanJobs(t *testing.T) {
	defer dontPanic(t)
	que := New(50)
	executeAssyncTest(t, que, 50, 10)
}

func TestBig(t *testing.T) {
	defer dontPanic(t)
	que := New(500)
	executeAssyncTest(t, que, 500, 3000)
}

func TestMult(t *testing.T) {
	defer dontPanic(t)

	que := NewMult(Mult{
		"a": 1,
		"b": 5,
		"c": 10,
		"d": 100,
		"e": 1000,
	})

	executeAssyncTestWithMult(t, que, "a", 1, 100)
	executeAssyncTestWithMult(t, que, "b", 5, 100)
	executeAssyncTestWithMult(t, que, "c", 10, 100)
	executeAssyncTestWithMult(t, que, "d", 100, 100)
	executeAssyncTestWithMult(t, que, "e", 1000, 100)
}
