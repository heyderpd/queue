## About Queue - simple use and control of queues

  This library was born, with the intention of solving the problem of starvation of http requests.
  Done in an application, that needed to parallel many different tasks depending on other micro-services.
  Is a simple solution that allows the control of parallel requests of any means.
  These need to limit the amount of concurrent active processes, originating from a large stack of tasks.

## Installing

  $ go get github.com/heyderpd/queue

## Simple Example

```go
  import (
    "github.com/heyderpd/queue"
  )

  var (
    // init your pool of queues with your limit
    que = queue.New(100)
    var jobsToDo = make([]string, 10000)
  )

  func main() {
    for _, name := range jobsToDo {
      go func(){
        // get next queue control
        mutex := que.Get()
        // use a simple mutex
        mutex.Lock()
        defer mutex.Unlock()

        /* and do your async things */
        time.Sleep(time.Millisecond * time.Duration(100))
      }
    }

    /* don't forget to control your mult async result */
  }
```

## Http Concurrent Example

The first 300 tasks were put into execution, all the rest were to wait for the completion of the first ones.
And will be evenly distributed among the queues. The only you need is Init() and Get(), and letting mutex logic do the magic!

```go
  var (
    // limit the amount of concurrent active processes
    que = queue.New(300)
    jobsToDo = make([]string, 99999)
  )

  func doHttpRequest(job string) {
    mutex := que.Get()
    /* call mutex lock and defer unlock in your function */
    mutex.Lock()
    defer mutex.Unlock()

    /* make your http request */
    client := &http.Client{}
    req, err := http.NewRequest(method, url + job, body)
    resp, err := client.Do(req)
    /* look for a coffee */
  }

  func main() {
    for _, job := range jobsToDo {
      /* start all requisition in the same time!!! =D */
      go doHttpRequest(job)
    }
  }
```

## Http Concurrent Example with Mult Queue

If you need balance concorrent request with diferent queues, with queue.NewMult you can create a map of queue groups!
Only change is the init need receive a specific map and the Get function substituted to que.GetGroup, that recei the group name to work.

```go
  var (
    // limit the amount of concurrent active processes by group!
    que = queue.NewMult(queue.Mult{
      "get-from-a": 10,
      "get-from-b": 20,
    })
    jobsToDo = make([]string, 1234)
  )

  func doHttpRequest(job string) {
    /* get mutex from specific queue group */
    mutex := que.GetGroup(job)

    /* call mutex lock and defer unlock in your function */
    mutex.Lock()
    defer mutex.Unlock()

    /* make your http request */
    client := &http.Client{}
    req, err := http.NewRequest(method, url + job, body)
    resp, err := client.Do(req)
    /* look for a coffee */
  }

  func main() {
    for _, job := range jobsToDo {
      /* start all requisition in the same time!!! =D */
      go doHttpRequest("get-from-a")
      go doHttpRequest("get-from-b")
    }
  }
```
