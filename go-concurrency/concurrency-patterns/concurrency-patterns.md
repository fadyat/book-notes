# Chapter 4

## Concurrency Patterns in Go

We've explored concurrency primitives, this chapter will deep dive
how to compose them into patterns that will help keep your system
scalable and maintainable.

> Most of the patterns will use `interface{}` to allow for any type
> of data to be passed around.
>
> To avoid them, you can use generics or generators.

## Available Patterns

- [Confinement](#confinement)
- [The `for-select` loop](#the-for-select-loop)
- [Preventing Goroutine Leaks](#preventing-goroutine-leaks)
- [The `or-channel`](#the-or-channel)
- [Error Handling](#error-handling)
- [Pipelines](#pipelines)
    - [Best Practices for Constructing Pipelines](#best-practices-for-constructing-pipelines)
    - [Some Handy Generators](#some-handy-generators)
- [Fan-Out, Fan-In](#fan-out-fan-in)
- [The `or-done` channel](#the-or-done-channel)
- [The `tee`-channel](#the-tee-channel)
- [The `bridge`-channel](#the-bridge-channel)
- [Queuing](#queuing)
- [The `context` package](#the-context-package)

### Confinement

Confinement is the simple yet powerful idea of ensuring that data
is only ever accessed by _one_ concurrent process.

There are two types of confinement:

- ad-hoc
  > when you achieve confinement through a convention - whether it be set by the
  > language community or your own team.

  ```go
  data := make([]int, 4)
  
  loopData := func(handleData chan<- int) {
      defer close(handleData)
      for i := range data {
        handleData <- data[i]
      }
  }
  
  handleData := make(chan int)

  // by the convention of our group, only loopData can access data
  go loopData(handleData)

  for num := range handleData {
      fmt.Println(num)
  }
  ```

  > code is touched by many people and can be hard to enforce

- lexical
  > when you achieve confinement through the lexical structure of your code.

  ```go
  chanOwner := func() <-chan int {
      results := make(chan int, 5)
      
      go func() {
          defer close(results)
          
          // do some work
      }()
      
      return results
  }
  
  consumer := func(results <-chan int) {
      for result := range results {
          fmt.Printf("Received: %d\n", result)
      }
      
      fmt.Println("Done receiving!")
  }
  
  results := chanOwner()
  consumer(results)
  ```

  > language enforces the confinement

It can be difficult to establish confinement, and so sometimes we have to fall back to our wonderful
Go concurrency primitives.

### The `for-select` loop

Cases when you can use it:

- Sending iteration variables out on a channel

  ```go
  for _, s := range []string{"a", "b", "c"} {
      select {
      case <-done:
          return
      case stringStream <- s:
      }
  }
  ```

- Looping infinitely waiting to be stopped

  ```go
  for {
      select {
      case <-done:
          return
      default:
      }
      
      // do some work, or do it in the `default` case
  }
  ```

### Preventing Goroutine Leaks

Goroutines are not garbage collected, so if you don't clean them up, they may
continue to run if you made a mistake.

Goroutine has a few paths to termination:

- when it completed it work
- when it couldn't continue its work due to a panic
- when it's told to stop working

Goroutine leak example:

```go
doWork := func(strings <-chan string) <-chan interface{} {
    completed := make(chan interface{})
    
    go func() {
        defer fmt.Println("doWork exited.")
        defer close(completed)
        
        for s := range strings {
            // do something interesting
            fmt.Println(s)
        }
    }()
    
    return completed
}

// reading from a `nil` channel will block forever
doWork(nil)

// long running program
```

The way to successfully clean up a goroutine is to use a `done` channel:

```go
doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
    completed := make(chan interface{})
    
    go func() {
        defer fmt.Println("doWork exited.")
        defer close(completed)
        
        for {
            select {
            case s := <-strings:
                // do something interesting
                fmt.Println(s)
            case <-done:
                return
            }
        }
    }()
    
    return completed
}

done := make(chan interface{})
completed := doWork(done, nil)

go func() {
    // cancel the operation after 1 second
    time.Sleep(1 * time.Second)
    fmt.Println("Canceling doWork goroutine...")
    close(done)
}()

<-completed
fmt.Println("Done.")
```

> The `done` channel is a common pattern in Go, and is used to signal
> cancellation or completion of a goroutine.

Case when goroutine is blocked to write a value to a channel:

```go
newRandStream := func() <-chan int {
    stream := make(chan int)
    
    go func() {
        defer fmt.Println("closure exited.")
        defer close(stream)
        
        for {
            stream <- rand.Int()
        }
    }()
    
    return stream
}

randStream := newRandStream()

fmt.Println("3 random ints:")
for i := 1; i <= 3; i++ {
    fmt.Printf("%d: %d\n", i, <-randStream)
}
```

> The goroutine is blocked on the `stream <- rand.Int()` line, and so
> the `defer` statements are never executed.

To fix this we again can use a `done` channel:

```go
newRandStream := func(done <-chan interface{}) <-chan int {
    stream := make(chan int)
    
    go func() {
        defer fmt.Println("closure exited.")
        defer close(stream)
        
        for {
            select {
            case stream <- rand.Int():
            case <-done:
                return
            }
        }
    }()
    
    return stream
}

done := make(chan interface{})
randStream := newRandStream(done)

fmt.Println("3 random ints:")
for i := 1; i <= 3; i++ {
    fmt.Printf("%d: %d\n", i, <-randStream)
}

close(done)
```

### The `or-channel`

At times, you may find yourself wanting to combine one or more done channels into a single done channel that closes if
any of its component channels close.

Sometimes you can't know the number of done channels you're working with a runtime.

This pattern is useful to employ at the intersection of modules in your system. At these intersections, you tend to have
multiple conditions for canceling trees of goroutines through your call stack. Using the or function, you can simply
combine these together and pass it down the stack. (Can be done with `context` package)

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// or function accepts a variadic number of channels and returns a single channel.
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		// this is recursive function -> stop point
		case 0:
			return nil
		// have only one channel -> return it
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})

		// we can wait for messages from multiple channels w/o blocking
		go func() {
			defer close(orDone)

			switch len(channels) {
			// an optimization, every recursive call will have at least two channels
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			// recursively create an or-channel from all channels and then select from it
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})

		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// we will see that the function returns after 1 second
	fmt.Printf("done after %v", time.Since(start))
}
```

### Error Handling

In concurrent programs, error handling can be difficult.
We often think about success cases first, and then we think about error cases.

The most fundamental question when thinking about error handling is, "Who should be responsible for handling the error?"

Separate your concerns: in general, your concurrent processes should send their errors to another part of your program
that has complete information about the state of your program, and can make a more informed decision about what to do.

```go
package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)

		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}

				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			continue
		}

		fmt.Printf("response: %v\n", result.Response.Status)
	}
}
```

If your goroutine can produce errors, those errors should be tightly coupled with your result type and passed along
through the same lines of communication—just like regular synchronous functions.

### Pipelines

A pipeline is just another tool you can use to form an abstraction in your system. In particular, it is a very powerful
tool to use when your program needs to process streams, or batches of data.

A pipeline is nothing more than a series of things that take data in, perform an operation on it, and pass the data back
out. We call each of these operations a stage of the pipeline.

By using a pipeline, you separate the concerns of each stage, which provides numerous benefits.

You can modify stages independent of one another, you can mix and match how stages are combined independent of modifying
the stages, you can process each stage concurrent to upstream or downstream stages, and you can _fan-out_, or
_rate-limit_ portions of your pipeline.

```go
package main

import "fmt"

func main() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}

		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}

		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}
```

Stage properties:

- stage consumes and returns the same type
- stage must be reified by the language so that it may be passed around.

Notice how each stage is taking a slice of data and returning a slice of data?

These stages are performing what we call **batch processing**.
> This just means that they operate on chunks of data all at once instead of one discrete value at a time.

There is another type of pipeline stage that performs **stream processing**.
> This means that the stage receives and emits one element at a time.

#### Best Practices for Constructing Pipelines

Channels are uniquely suited to constructing pipelines in Go because they fulfill all of our basic requirements.

```go
package main

import "fmt"

func main() {

	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		stream := make(chan int)

		go func() {
			defer close(stream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case stream <- i:
				}
			}
		}()

		return stream
	}

	multiply := func(
		done <-chan interface{},
		stream <-chan int,
		multiplier int,
	) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			defer close(multipliedStream)

			for i := range stream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	add := func(
		done <-chan interface{},
		stream <-chan int,
		additive int,
	) <-chan int {
		addedStream := make(chan int)

		go func() {
			defer close(addedStream)

			for i := range stream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()

		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
```

Each stage executes concurrently with the others, and each stage is decoupled from the others.

Our entire pipeline is always preemptable by closing the `done` channel.

#### Some Handy Generators

```go
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// will generate the same values until not stopped
	repeat := func(
		done <-chan interface{},
		values ...interface{},
	) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case stream <- v:
					}
				}
			}
		}()

		return stream
	}

	// will take the first num values from the stream
	take := func(
		done <-chan interface{},
		stream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-stream:
				}
			}
		}()

		return takeStream
	}

	// will generate values from fn until not stopped
	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				select {
				case <-done:
					return
				case stream <- fn():
				}
			}
		}()

		return stream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1, 2, 3), 10) {
		fmt.Println(num)
	}

	rnd := func() interface{} { return rand.Int() }
	for num := range take(done, repeatFn(done, rnd), 10) {
		fmt.Println(num)
	}

	// will convert the stream to chan string
	toString := func(
		done <-chan interface{},
		stream <-chan interface{},
	) <-chan string {
		stringStream := make(chan string)

		go func() {
			defer close(stringStream)

			for v := range stream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()

		return stringStream
	}

	for str := range toString(done, take(done, repeat(done, "a", "b", "c"), 10)) {
		fmt.Println(str)
	}
}
```

Type-specific stages are twice as fast, but only marginally faster in magnitude. Generally, the limiting factor on your
pipeline will either be your generator, or one of the stages that is computationally intensive.

### Fan-Out, Fan-In

So you’ve got a pipeline set up. Data is flowing through your system beautifully, transforming as it makes its way
through the stages you’ve chained together. It’s like a beautiful stream; a beautiful, slow stream, and oh my god why is
this taking so long?

Sometimes, stages in your pipeline can be particularly computationally expensive.
When this happens, upstream stages in your pipeline can become blocked while waiting for your expensive stages to
complete. Not only that, but the pipeline itself can take a long time to execute as a whole. How can we address this?

**Fan-out** is a term to describe the process of starting multiple goroutines to handle input from the pipeline, and
**fan-in** is a term to describe the process of combining multiple results into one channel.

You might consider fanning out one of your stages if both of the following apply:

- it doesn't rely on values that the stage had calculated previously
- it takes a long time to complete

The property of order-independent is important because we don't have guarantee
in what order concurrent copies of your stage will run, nor in what order they will return.

```go
package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	// fanIn is a function that takes a variadic number of channels
	// 	and multiplexes them onto a single channel.
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) {
			defer wg.Done()

			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
	) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				select {
				case <-done:
					return
				case stream <- fn():
				}
			}
		}()

		return stream
	}

	take := func(
		done <-chan interface{},
		stream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-stream:
				}
			}
		}()

		return takeStream
	}

	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}

		for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
			if n%i == 0 {
				return false
			}
		}

		return true
	}

	primeFinder := func(
		done <-chan interface{},
		stream <-chan interface{},
	) <-chan interface{} {
		primeStream := make(chan interface{})

		go func() {
			defer close(primeStream)

			for {
				select {
				case <-done:
					return
				case i := <-stream:
					if isPrime(i.(int)) {
						primeStream <- i
					}
				}
			}
		}()

		return primeStream
	}

	done := make(chan interface{})
	defer close(done)

	rnd := func() interface{} { return rand.Intn(5_000_000) }
	stream := repeatFn(done, rnd)

	findersNumber := runtime.NumCPU()
	finders := make([]<-chan interface{}, findersNumber)
	fmt.Printf("Spinning up %d prime finders.\n", findersNumber)
	for i := 0; i < findersNumber; i++ {
		finders[i] = primeFinder(done, stream)
	}

	start := time.Now()
	fmt.Println("Primes:")
	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))
}
```

#### The `or-done` channel

That is to say, you don’t know if the fact that your goroutine was canceled means the channel you’re reading from will
have been canceled. We need to wrap our read from the channel with a select statement that also selects from a done
channel. This is perfectly fine, but doing so takes code that’s easily read like this:

```go
package main

import "fmt"

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if !ok {
						return
					}

					select {
					case <-done:
					case stream <- v:
					}
				}
			}
		}()

		return stream
	}

	done := make(chan interface{})
	defer close(done)

	values := make(chan interface{})

	go func() {
		defer close(values)
		for i := 0; i < 10; i++ {
			values <- i
		}
	}()

	for val := range orDone(done, values) {
		fmt.Println(val)
	}
}
```

### The `tee`-channel

Sometimes you may want to split values coming in from a channel so that you
can send them off into two separate areas of your codebase.

Taking its name from the `tee` command in Unix, which reads from standard input
and writes to standard output and files, the `tee` channel is a channel that
splits incoming data into two copies of the data.

```go
package main

import "fmt"

func main() {
	var repeat func(done <-chan interface{}, values ...interface{}) <-chan interface{}
	var take func(done <-chan interface{}, stream <-chan interface{}, num int) <-chan interface{}
	var orDone func(done, c <-chan interface{}) <-chan interface{}

	tee := func(
		done <-chan interface{},
		in <-chan interface{},
	) (_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})

		go func() {
			defer func() {
				close(out1)
				close(out2)
			}()

			for val := range orDone(done, in) {
				// we want to use local versions of out1 and out2,
				//  so we shadow them
				var out1, out2 = out1, out2

				// making for loop to send val to out1 and out2
				//  don't block each other
				for i := 0; i < 2; i++ {
					// once we read from out1 or out2, we set shadowed
					//  copies to nil, so that we don't send to nil channels
					select {
					case <-done:
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()

		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2, 3, 4, 5), 10))
	for val1 := range out1 {
		fmt.Printf("out1: %v\n", val1)
		fmt.Printf("out2: %v\n", <-out2)
	}
}
```

Notice that writes to `out1` and `out2` are tightly coupled. The iteration over in cannot continue until both `out1`
and `out2` have been written to. Usually this is not a problem as handling the throughput of the process reading from
each channel should be a concern of something other than `the` tee command anyway, but it’s worth noting.

### The `bridge`-channel

In some circumstances, you may find yourself wanting to consume values from a sequence of channels:

```go
<-chan <-chan interface{}
```

As a consumer, the code may not care about the fact that its values come from a sequence of channels. In that case,
dealing with a channel of channels can be cumbersome. If we instead define a function that can destructure the channel
of channels into a simple channel - a technique called `bridging` the channels - this will make it much easier for the
consumer to focus on the problem at hand.

```go
package main

import "fmt"

func main() {
	var orDone func(done <-chan interface{}, c <-chan interface{}) <-chan interface{}

	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		// this channel will return all the values from the channels
		stream := make(chan interface{})

		go func() {
			defer close(stream)

			for {
				// pooling the channels
				var ch <-chan interface{}
				select {
				case maybeCh, ok := <-chanStream:
					if !ok {
						return
					}

					ch = maybeCh
				case <-done:
					return
				}

				// reading the values from the channel
				for val := range orDone(done, ch) {
					select {
					case stream <- val:
					case <-done:
					}
				}
			}
		}()

		return stream
	}

	genValues := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))

		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()

		return chanStream
	}

	done := make(chan interface{})
	defer close(done)

	for val := range bridge(done, genValues()) {
		fmt.Println(val)
	}
}
```

Thanks to bridge, we can use the channel of channels from within a single range statement and focus on our loop’s logic.
Destructuring the channel of channels is left to code that is specific to this concern.

### Queuing

Sometimes it’s useful to begin accepting work for your pipeline even though the pipe‐ line is not yet ready for more.
This process is called **queuing**.

All this means is that once your stage has completed some work, it stores it in a temporary location in memory so that
other stages can retrieve it later, and your stage doesn't need to hold a reference to it. (buffered channels are a type
of queue)

While introducing queuing into your system is very useful, it’s usually one of the last techniques you want to employ
when optimizing your program. Adding queuing prematurely can hide synchronization issues such as deadlocks and
livelocks, and further, as your program converges toward correctness, you may find that you need more or less queuing.

So the answer to our question of the utility of introducing a queue isn’t that the runtime of one of stages has been
reduced, but rather that the time it’s in a **blocking state** is reduced.

In this way, the true utility of queues is to decouple stages so that the runtime of one stage has no impact on the
runtime of another. Decoupling stages in this manner then cascades to alter the runtime behavior of the system as a
whole, which can be either good or bad depending on your system.

Let’s begin by analyzing situations in which queuing can increase the overall performance of your system. The only
applicable situations are:

- if batching requests in a stage saves time
- if delays in a stage produce a feedback loop into the system

Queueing should be implemented either:

- at the entrance to your pipeline
- in stages where batching will lead to higher efficiency

Queuing can be useful in your system, but because of its complexity, it’s usually one of the last optimizations which
should be applied.

### The `context` package

In concurrent programs it’s often necessary to preempt operations because of timeouts, cancellation, or failure of
another portion of the system.

It turns out that the need to wrap a done channel with this information is very common in systems of any size, and so
the Go authors decided to create a standard pattern for doing so.

- `Done()` returns a channel that’s closed when the context is canceled or times out.
- `Err()` returns a non-nil error value after the context is canceled or times out.
- `Deadline()` returns the time at which the context will time out.

Context serves two purposes:

- to provide API for cancelling branches of your call graph
- to provide big data-bag for transporting request-scoped values through your call graph

Cancellation is a fn of three aspects:

- goroutine's parent may want to cancel it
- goroutine may want to cancel its children
- any blocking operations within a goroutine need to be preemptable so that it may be cancelled

If you look at the methods on the Context interface, you'll see that there's nothing present that can mutate the state
of the underlying structure. Further, there’s nothing that allows the function accepting the Context to cancel it. This
protects functions up the call stack from children canceling the context. Combined with the Done method, which provides
a done channel, this allows the Context type to safely manage cancellation from its antecedents.

By using `context.WithCancel`, `context.WithDeadline`, and `context.WithTimeout`, you can add a cancelable context to
your call graph.

At the top of your asynchronous call-graph, your code probably won’t have been passed a Context. To start the chain, the
context package provides you with two functions to create empty instances of Context: `context.Background`
and `context.TODO`.

> `Background` simply returns an empty Context.
>
> `TODO` is not meant for use in production, but also returns an empty Context;
>
> `TODO`’s intended purpose is to serve as a placeholder for when you don't know which Context to utilize, or if
> you expect your code to be provided with a Context, but the upstream code hasn't yet furnished one.

Since both the Context’s key and value are defined as interface{}, we lose Go’s type-safety when attempting to retrieve
values. The key could be a different type, or slightly different from the key we provide. The value could be a different
type than we’re expecting. For these reasons, the Go authors recommend you follow a few rules when storing and
retrieving value from a Context:

- define a custom key-type for your package
  > Since the type you define for your package’s keys is unexported, other packages cannot conflict with keys you
  > generate within your package.

- since we don’t export the keys we use to store the data, we must therefore export functions that retrieve the data for
  us.

```go
type ctxKey int

const (
    ctxUserID ctxKey = iota
    ctxAuthToken
)

func UserID(c context.Context) string {
    return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
    return c.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
    ctx := context.WithValue(context.Background(), ctxUserID, userID)
    ctx = context.WithValue(ctx, ctxAuthToken, authToken)

    HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
    userID := UserID(ctx)
    authToken := AuthToken(ctx)

    // do something with userID and authToken
}
```

The context package is pretty neat, but it hasn't been uniformly lauded. Within the Go community, the context package
has been somewhat controversial. The cancellation aspect of the package has been pretty well received. Still, the
ability to store arbitrary data in a Context, and the type-unsafe manner in which the data is stored, have caused some
divisiveness. Although we have partially abated the lack of type safety with our accessor functions, we could still
introduce bugs by storing incorrect types. However, the more significant issue is definitely the nature of what
developers should store in instances of Context.

The most prevalent guidance on what’s appropriate is this somewhat ambiguous comment in the context package:

> Use context values only for request-scoped data that transits processes and
> API boundaries, not for passing optional parameters to functions.

What is "request-scoped data"?

Here some author heuristics:

- the data should transit process and API boundaries
- the data should be immutable
- the data should be trand toward simple types
- the data should be data, not types with methods
- the data should help decorate operations, not drive them

Another dimension to consider is how many layers this data might need to traverse before utilization.

The cancellation functionality provided by Context is very useful, and your feelings about the data-bag shouldn't deter
you from using it.