### Scheduler

OS threads are slow.

Golang creators decided to use goroutines (green threads) instead of OS threads.

- Goroutines are cheap
- Goroutines are not OS threads
- Faster context switching

Goroutines are multiplexed to fewer number of OS threads

![](./docs/mpg.jpeg)

#### Structure

- P: Processor
- M: Machine
- G: Goroutine 

P = number of processors, often number of cores
> logical processors, number are configurable using `runtime.GOMAXPROCS(n int) int`
> 
> on each logical processors can be executed only one goroutine at a time 

M = number of OS threads, often equal to P
> each P works on one M thread. 
>
> P != M for example if M is blocked by syscall, also have CGO and other cases
>
> each P have a self managed work queue of Gs, which she executes one by one.

P, M, G are not real OS objects, they are just structures in Go runtime.

![](./docs/mpg-struct.jpeg)

Why we need to separate M, P but they are always nested?
Because os threads can be blocked, then scheduler unbinds G from M, P and binds it to another.

### Queues 

- GRQ(Global Run Queue)
- LRQ(Local Run Queue)
> each P have a LRQ, which contains Gs, which are ready to execute 

#### Multitasking types

- Preemptive multitasking 
> All tasks are equal. For each task equal time is allocated.
>
> Example: OS scheduler 
>
> Based on some metrics reassign resources to another threads.

- Cooperative multitasking 
> Tasks executes as long as they want. Tasks should yield control to other tasks.
>
> Sleeps until one of the tasks aren't woke up him with a request to yield control.
> Then scheduler decides which task should be executed next.

Golang uses some kind of cooperative multitasking.

- Goroutines allow others goroutines to run, when they are calling blocking functions (IO, channel, OS calls etc)
- Goroutines allow others goroutines to run, when they are calling `runtime.Gosched()`
- Goroutines allow others goroutines to run, when making function calls.

#### Goroutine states

- Waiting
> is stopped and waiting for something in order to continue.
>
> Reasons like waiting for IO, sync calls, OS etc.

- Runnable
> wants time on an M so it can execute its assigned instructions. 

- Executing
> has been placed on an M and is executing its instructions.

#### Core concepts

- FIFO queue of Gs
> don't have priorities, all Gs are equal

- Create minimal number of threads
- When threads are free, they will be reused

- No one can interrupt G 
> no execution time guarantees

`go func == runtime.newproc(func)`

> When we call `go func` we create a new G and put it to the end of the queue, not executed immediately.

Executed using `func schedule()` in `runtime/proc.go`

Short representation of `schedule()`:

```go
func schedule() {
    // are we need a GC?
    g := get_gc_worker()
    if g == nil {
        // get G from queue
        g = runqget()
    }

    if g == nil {
        // steal G from other M or wait until G will be available
        g = steal_or_wait()
    }

    // execute G
    execute(g)
}
```

Who calls `schedule()`? - `runtime.mstart()` from assembler code `runtime/asm_amd64.s`

Then called `main` in `runtime/proc.go`

```go
func main() {
    go start_gc()
    // ... 
    main_main() // actual main function
}
```

`schedule()` also called in `goexit()`.

Repeat:

- When starting program, in assembler code called first round of `schedule()`
- Scheduler execures `runtime/proc.go:main()`
- `runtime/proc.go:main()` calls user `main()`
- when goroutine is finished, `goexit()` is called with `schedule()`

> not recursive realisation

Scheduler is not separate thread, just called between goroutines.
For each M separate scheduler.

#### Sources

- [Go under the hood, RU](https://www.youtube.com/watch?v=rloqQY9CT8I&t=1683s)
- [Scheduling in Go](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)
