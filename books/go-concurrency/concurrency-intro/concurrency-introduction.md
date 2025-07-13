# Chapter 1

## An introduction to concurrency

Concurrency, asynchronous programming, parallelism and threading are all terms that are often used interchangeably, but
they are not the same thing.

Practical example of concurrent model are people. You are currently reading this sentence, while others in the world are
simultaneously living their lives.

### Moores's law, Web scale, and the mess we're in

Moore's law is a prediction that the number of transistors in a dense integrated circuit will double approximately every
two years.

Multicore processors were created, because of slowdown in the Moore's law. This looks like a clever solution, but
computer scientists soon found themselves facing down the limits of another law: Amdahl's law.

Amdahl's law is a formula that describes the maximum speedup that can be achieved by implementing solutions in a
parallel manner. Simply put, it states that the maximum speedup is limited by the portion of the program that cannot be
parallelized.

For example, we have a GUI based program. It depends on from the user input, and it is not possible to parallelize it.

And we have a program, that computes the pi number. It's possible to parallelize it.

**Embarrassingly parallel** -- a problem that can be easily divided into parallel tasks.

Amdahl's law helps us to understand the difference between this two problems, and can help us decide whether
parallelization is the right solution for our problem.

**Horizontal scaling** -- adding more machines or CPUs to a system to improve performance.

**Vertical scaling** -- adding more memory or CPU power to a single machine to improve performance.

**Web scale** -- a term used to describe the need for a system to be able to handle a large number of users.

### Why is concurrency hard?

#### Race conditions

**Race condition** occurs when two or more operations must be executed in a specific order, but the order is not
guaranteed.

Most of the time, this show in **data race**, where one concurrent operation is trying to read data, while another
concurrent operation is trying to write the same data.

#### Atomicity

**Atomicity** -- property of an operation that guarantees that it will be executed as a single unit, without any
interruptions.

Atomic operations are implicitly safe within concurrent context.

Most operations are not atomic, and we need to use special tools to make them atomic.

#### Memory access synchronization

**Critical section** -- a part of program, where we need exclusive access to shared data.

Memory access synchronization is simply demonstrates with the following example:

It's not _go_ idiomatic, but good for understanding.

```go
var x int64
var wg sync.Mutex
go func() {
    wg.Lock()
    x++
    wg.Unlock()
}()

wg.Lock()
if x == 0 {
    fmt.Println("x is zero")
} else {
    fmt.Println("x is not zero")
}
```

We have solved data race, but not race condition. The order of executions a still non-deterministic.

You can solve some problems by synchronizing access to the memory, but as we just saw, it doesn’t
automatically solve data races or logical correctness. Further, it can also create maintenance and performance problems.

#### Deadlocks

**Deadlock** -- a situation where two or more concurrent operations are blocked forever, waiting for each other to
complete.

```go
type Value struct {
    mu sync.Mutex
    value int
}

var wg sync.WaitGroup
printSum := func(v1, v2 *Value) {
    defer wg.Done()
    v1.mu.Lock()
    defer v1.mu.Unlock()
    time.Sleep(2 * time.Second)
    v2.mu.Lock()
    defer v2.mu.Unlock()
    fmt.Printf("sum=%v\n", v1.value + v2.value)
}

var a, b Value
wg.Add(2)
go printSum(&a, &b)
go printSum(&b, &a)
wg.Wait()
```

Coffman's conditions are the basic techniques that help detect, prevent and correct deadlocks.

- Mutual exclusion -- only one process can access a resource at a time.
- Wait for condition -- a process may be allocated some resources while waiting for others.
- No preemption -- a resource cannot be taken away from a concurrent process until it releases the resource.
- Circular wait -- a set of concurrent processes are waiting for each other in a circular chain.

These conditions are necessary, but not sufficient for deadlock to occur.

#### Livelocks

**Livelock** -- a situation that are actively performing concurrent operations, but are unable to make progress.

Real world example:

- Two people are trying to pass each other on a narrow sidewalk. They are both moving, but neither is able to
  complete the pass.

Example: [livelock.go](https://github.com/GermanGorelkin/go-patterns/blob/master/concurrency/problems/deadlocks-livelocks-and-starvation/livelock/main.go)

Livelocks are a subset of a larger problems called **starvation**.

#### Starvation

**Starvation** -- a situation where a concurrent operation can't get all the resources it needs to perform work.

More broadly, starvation usually implies that there are one or more greedy concur‐ rent process that are unfairly
preventing one or more concurrent processes from accomplishing work as efficiently as possible, or maybe at all.

Example: [starvation.go](https://github.com/GermanGorelkin/go-patterns/blob/master/concurrency/problems/deadlocks-livelocks-and-starvation/starvation/main.go)

### Determining concurrency safety

If you’re starting with a blank slate and need to build up a sensible way to model your problem space and
concurrency is involved, it can be difficult to find the right level of abstraction.

Not always obvious, how to model the problem.

### Simplicity in the face of complexity

Concurrency is certainly a difficult area in CS, but these problems are not intractable.

Memory management can be another difficult problem domain in computer science, and when combined with concurrency,
it can become extraordinarily difficult to write correct code. Go made memory management easier by introducing
garbage collection.

Go runtime also automatically handles multiplexing concurrent operations onto OS threads. 
It allows you to directly map concurrent problems into concurrent constructs instead of dealing with the
minutia of starting and managing threads, and mapping logic evenly across available threads.

Go concurrency primitives also make composing larger problems easier. Channels are composable, concurrent-safe way to
communicate between concurrent processes.

