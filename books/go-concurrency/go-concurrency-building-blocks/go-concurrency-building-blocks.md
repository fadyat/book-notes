# Chapter 3

## Go's Concurrency Building Blocks

### Goroutines

Go program have at least one goroutine, the main goroutine, which is created automatically when the program begins.

You can create new goroutines using the `go` keyword followed by a function invocation.

```go
go func() {
    // do something
}()
```

Works with predefined functions, lambdas, anonymous functions, and methods.

Goroutines are not OS threads - threads are managed by a language runtime -
there're higher level of abstraction known as **coroutines**.

**coroutines** are simple concurrent subroutines(functions, clojures, methods)
that are non-preemptive (they can't be interrupted), instead they have
multiple points throughout which allow for suspension or reentry.

Goroutines are deep integreated with the Go runtime, he observes behavior
of goroutines and automatically suspends them when they block and then resumes them when they become unblocked.

Go’s mechanism for hosting goroutines is an implementation of what’s called
a **M:N scheduler**, which means it maps M green threads to N OS threads.
Goroutines are then scheduled onto the green threads by the Go runtime scheduler.
When we have more goroutines than green threads available, the scheduler handles the distribution of the goroutines
across the available threads and ensures that when these goroutines become blocked, other goroutines can be run.

Go follows a model of concurrency called the fork-join model.
> fork = in any point of time, program can split off a child branch of execution
> to run concurrent with the parent branch
>
> join = at some point in the future, the child branch of execution can rejoin the parent branch

We can use the `sync.WaitGroup` type to synchronize the execution of goroutines.

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // do something
}()

wg.Wait()
```

Goroutines are executed in the same address space, where they were created,
when passing value from loop to goroutine, we need to pass a copy of the value,
because the value can be changed in the loop before the goroutine is executed.

```go
for _, val := range values {
    go func(val int) {
        defer wg.Done()
        // do something with val
    }(val)
}
```

If goroutines are still using some variable which is out of scope,
the variable will be moved to the heap, so the goroutines can still access it.

GC does nothing with goroutines, which been abandoned some how. They will
hang around until the program exits.

Goroutines are lightweight, they only take 2KB of memory, so we can create
millions of them.

**Context switching** is when something hosting a concurrent process must save
its state to switch to running a different concurrent process.

If we have many concurrent processes, we can spend all of our CPU time
context switching between them, instead of actually doing work.

At the OS level is costly to context switch between threads, because the OS
must save and restore registers, lookup tables, memory mappings, and other
resources.

Context switching in software is comparatively much, much cheaper. Under a software
defined scheduler, the runtime can be more selective in what is persisted
for retrieval, how it is persisted, and when the persisting need occurs.

### The `sync` Package

The `sync` package provides low-level concurrency primitives, such as mutexes,
condition variables, wait groups, and so on.

#### WaitGroup

`sync.WaitGroup` is a great way to wait for a set of concurrent operations to
complete when you either don't care about the result of the concurrent operation,
or you have other means of collecting the result.

```go
var wg sync.WaitGroup

wg.Add(1)
go func() {
    defer wg.Done()
    // do something
}()

wg.Add(1)
go func() {
    defer wg.Done()
    // do something
}()

wg.Wait()
```

`WaitGroup` = concurrent-safe counter, `Add` increments the counter, `Done`
decrements the counter, `Wait` blocks until the counter is zero.

#### Mutex and RWMutex

`sync.Mutex` = mutual exclusion lock, it's a way to guard a critical section of
your code.

Critical section = a section of code where program require exclusive access to
some shared resource.

```go
var mu sync.Mutex 
var balance int 

func deposit(amount int) {
    mu.Lock()
    defer mu.Unlock()
    balance += amount
}
```

`Unlock` must be called in a `defer` statement, so it will be called even if
the function panics.

Minimize the amount of work you do while holding a lock, because it can
decrease the performance of your program.

If not all processes are writing to the shared resource, we can use
`sync.RWMutex` = reader/writer mutual exclusion lock.

With `RWMutex` we can have multiple readers or one writer at the same time.

```go
var mu sync.RWMutex
var mp = make(map[string]string)

func lookup(key string) string {
    mu.RLock()
    v := mp[key]
    mu.RUnlock()
    return v
}

func set(key, value string) {
    mu.Lock()
    mp[key] = value
    mu.Unlock()
}
```

#### Cond

`sync.Cond` = condition variable, it's a way to signal between goroutines that
something has happened.

Note: `Wait` call doesn't just block, it suspends the goroutine that calls it,
so it can be resumed later.

When we call `Wait` unlock called on the mutex, and uppon when exiting `Wait`
lock called on the mutex.

```go
package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// create a new condition variable
	c := sync.NewCond(&sync.Mutex{})

	// queue is empty at this point, later we add some entries
	queue := make([]interface{}, 0, 10)

	// wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	remove := func(delay time.Duration) {
		defer wg.Done()
		time.Sleep(delay)

		// entering the critical section, lock the mutex
		c.L.Lock()
		queue = queue[1:]
		log.Println("removed from queue: ", len(queue))

		// leaving the critical section, unlock the mutex
		c.L.Unlock()

		// signal the condition that the queue changed
		//
		// notifies goroutines, that are waiting for the condition
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)

		// entering the critical section, lock the mutex
		c.L.Lock()

		// waiting for the condition to be true
		for len(queue) == 2 {

			// suspends the execution of the calling goroutine
			c.Wait()
		}

		log.Println("adding to queue: ", len(queue))
		queue = append(queue, struct{}{})

		// calling new goroutine to remove from queue
		go remove(1 * time.Second)

		// leaving the critical section, unlock the mutex
		c.L.Unlock()
	}

	wg.Wait()
}
```

We can use `Broadcast` to notify all goroutines waiting for the condition.

```go
package main

import "sync"

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{
		sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()

			c.L.Lock()
			defer c.L.Unlock()

			c.Wait()
			fn()
		}()

		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		println("Maximizing window.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		println("Mouse clicked.")
		clickRegistered.Done()
	})

	// Broadcast wakes all goroutines waiting on c.
	button.Clicked.Broadcast()

	clickRegistered.Wait()
}
```

#### Once

`sync.Once` = a way to ensure that a function is called only once.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	increment := func() { count++ }

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()

	fmt.Println("Count is", count)
}
```

#### Pool

`sync.Pool` is a concurrent-safe implementation of the object pool pattern.

At a high level, a pool pattern is a way to create and make available a fixed
number of things to use.

It’s commonly used to constrain the creation of things that are expensive (e.g., database connections) so that only a
fixed number of them are ever created, but an indeterminate number of operations can still request access to these
things. In case of Go, it can be safely used concurrently.

So why use a pool and not just instantiate objects as you go? Go has a garbage collec‐ tor, so the instantiated objects
will be automatically cleaned up.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated++
			mem := make([]byte, 1024)
			return &mem // pointer to a slice
		},
	}

	// seeding the pool with 4KB
	for i := 0; i < 4; i++ {
		calcPool.Put(calcPool.New())
	}

	const numWorkers = 1024 * 1024

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte) // type assertion
			defer calcPool.Put(mem)

			// do something interesting with mem
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
```

Another common situation where a Pool is useful is for warming a cache of pre-allocated objects for operations that must
run as quickly as possible.

As we’ve seen, the object pool design pattern is best used either when you have con‐ current processes that require
objects, but dispose of them very rapidly after instantiation, or when construction of these objects could negatively
impact memory.

When working with a `sync.Pool`, it's important to:

- when instantiating the pool, give it a `New` member variable is thread-safe when called
- when you receive an instance from `Get`, make no assumptions regarding the state of the object you receive back.
- make sure to call `Put` when you're done with the object. (usually done with `defer`)
- objects in the pool must be roughly uniform

### Channels

Channels are one of the synchronization primitives in Go derived from Hoare’s Communicating Sequential Processes (CSP).

While they can be used to synchronize access to memory, they are best used to
communicate values between goroutines.

Can read and write to channels, channels are typed.

```go
stream := make(chan int)
```

Also we can define the direction of the channel. (corresponing to send and receive)

```go
receiveStream := make(<-chan int)
sendStream := make(chan<- int)
```

You will get a compile-time error if you try to send to a receive-only channel or receive from a send-only channel.

Channels in Go are **blocking**.

When reading from a channel can use a second variable to check if the channel is closed.

```go
val, ok := <-stream
```

In programs it's very useful to be able to indicate that no more values will be sent on a channel,
this can be done by closing the channel.

```go
stream := make(chan int)
close(stream)
```

When channel is closed we still can continue to read from it.

```go
package main

import (
	"fmt"
	"time"
)

func putSomeData(stream chan<- interface{}) {
	defer func() {
		close(stream)
		fmt.Println("close stream")
	}()

	for i := 0; i < 5; i++ {
		stream <- i
	}
}

func main() {
	// making buffered channel to close channel before reading
	stream := make(chan interface{}, 5)

	// blocking read = deadlock
	//  `ok` returned only if channel is closed or have some data
	//
	// v, ok := <-stream
	// fmt.Printf("%v, %v\n", v, ok)

	go putSomeData(stream)

	time.Sleep(time.Second)
	fmt.Println("sleep done")
	for i := 0; i < 10; i++ {
		v, ok := <-stream
		fmt.Printf("%v, %v\n", v, ok)
	}

	stream = make(chan interface{}, 5)
	go putSomeData(stream)

	// range will perform auto exit if channel is closed
	for v := range stream {
		fmt.Printf("%v\n", v)
	}
}
```

We can close channel to unblock all goroutines that are waiting on it. (like `sync.Cond`)

We can also create buffered channel - channel with a capacity.

```go
stream := make(chan int, 5)
```

Channels are blocking only when they are limited by their capacity.
Unbuffered channels are always blocking, because they have no capacity.

Default value for channel is `nil`, so if you try to send or receive from `nil`
channel you will get a runtime error.

Closing a `nil` or closed channel will also result in a runtime error.

#### Ownership.

The goroutine that owns the channel should:

- instantiate the channel
- perform writes or pass ownership to another goroutine
- close the channel
- encapsulate the previous three things in this list and expose them via a reader channel.

With such logic we are removing all possible `panic` from our code.

```go
package main

func main() {
	owner := func() <-chan int {
		ch := make(chan int, 5)

		go func() {
			defer close(ch)

			for i := 0; i < 5; i++ {
				ch <- i
			}
		}()

		return ch
	}

	results := owner()
	for result := range results {
		println("result:", result)
	}

	println("done")
}
```

#### The `select` Statement

The select statement allows to bind a set of send and receive operations together to wait for one of them to complete.

Getting first channel that is ready.

```go
var c1, c2 <-chan interface{}
var c3 chan<- interface{}

select {
case <-c1:
    // do something
case <-c2:
    // do something
case c3 <- struct{}{}:
    // do something
}
```

Selection of the channel is pseudo-random.

To avoid blocking we can use `time.After` channel.

```go
var c <-chan int

select {
case <-c:
case <-time.After(1 * time.Second):
    fmt.Println("Timed out.")
}
```

If we need to perform some action while waiting for a channel to be ready, we can use `default` case.
Usually you’ll see a default clause used in conjunction with a `for` loop.

```go
for {
    select {
    case <-c:
        // do something
    default:
        // skip
    }
    
    // do something while waiting
}
```

### The GOMAXPROCS Lever

Controls the number of OS threads that can execute Go code simultaneously.

Automatically set to the number of cores on the machine.


