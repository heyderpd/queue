package queue

import (
	"testing"
	"time"
	// "fmt"
)

var (
	sleepWeight = 3
	sleep = time.Millisecond * time.Duration(sleepWeight)
)

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

func job(que queuesControl, channel chan bool) {
	mutex := que.Get()
	mutex.Lock()
	defer mutex.Unlock()

	time.Sleep(sleep) // go some work
	channel <- true
}

func getNow() float64 {
	return float64(time.Now().UnixNano())
}

func getExpected(t int) float64 {
	return float64(t *1000 *1000 *sleepWeight) *1.15
}

func hasMuchTimeDIff(err float64, t *testing.T) {
	if 0.8 > err || err > 1.2 {
		t.Errorf("error: has so much time diff %v", err)
	}
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

	testQtd := 10
	expected := getExpected(testQtd)
	start := getNow()

	for range make([]bool, testQtd) {
		mutex := que.Get()
		mutex.Lock()
		time.Sleep(sleep) // go some work
		mutex.Unlock()
	}

	err := (getNow() - start) / expected
	hasMuchTimeDIff(err, t)
}

func TestSequencialWihtRoutines(t *testing.T) {
	defer dontPanic(t)

	que := New(1)

	channel := make(chan bool)
	testQtd := 10
	count := testQtd
	expected := getExpected(testQtd)
	start := getNow()

	go func(){
		time.Sleep(time.Millisecond * 1)

		for range make([]bool, testQtd) {
			go job(que, channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDIff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}

func TestWithTwoQueue(t *testing.T) {
	defer dontPanic(t)

	que := New(2)

	channel := make(chan bool)
	testQtd := 10
	count := testQtd
	expected := getExpected(testQtd) / 2
	start := getNow()

	go func(){
		time.Sleep(time.Millisecond * 1)

		for range make([]bool, testQtd) {
			go job(que, channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDIff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}

func TestWithFiveQueue(t *testing.T) {
	defer dontPanic(t)

	que := New(5)

	channel := make(chan bool)
	testQtd := 10
	count := testQtd
	expected := getExpected(testQtd) / 5
	start := getNow()

	go func(){
		time.Sleep(time.Millisecond * 1)

		for range make([]bool, testQtd) {
			go job(que, channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDIff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}

func TestWithTeenQueue(t *testing.T) {
	defer dontPanic(t)

	que := New(10)

	channel := make(chan bool)
	testQtd := 10
	count := testQtd
	expected := getExpected(testQtd) / 10 *1.1
	start := getNow()

	go func(){
		time.Sleep(time.Millisecond * 1)

		for range make([]bool, testQtd) {
			go job(que, channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDIff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}

func TestMoreQueueThanJobs(t *testing.T) {
	defer dontPanic(t)

	que := New(50)

	channel := make(chan bool)
	testQtd := 10
	count := testQtd
	expected := getExpected(testQtd) / 10 *1.1
	start := getNow()

	go func(){
		time.Sleep(time.Millisecond * 1)

		for range make([]bool, testQtd) {
			go job(que, channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDIff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}

func TestBig(t *testing.T) {
	defer dontPanic(t)

	que := New(500)

	channel := make(chan bool)
	testQtd := 3000
	count := testQtd
	expected := getExpected(testQtd) / 500
	start := getNow()

	go func(){
		time.Sleep(time.Millisecond * 1)

		for range make([]bool, testQtd) {
			go job(que, channel)
		}
	}()

	for {
		select {
		case <- channel:
			count--
			if count <= 0 {
				err := (getNow() - start) / expected
				hasMuchTimeDIff(err, t)
				return
			}
		case <- time.After(time.Second * 1):
			t.Errorf("end by timeout")
			return
		}
	}
}
