### Channels

- goroutine-safe (hchan mutex + some atomic)
- can save elements, FIFO semantics (hchan buffer)
- sharing data between goroutines (sendDirect, operations with buffer)
- blocking goroutines on reads/writes if the channel is empty/full (sendq/recvq, sudog, calls to scheduler: gopark(), goready())

All channels are stored in the heap.

```go
type hchan struct {
    qcount   uint           // total data in the queue
    dataqsiz uint           // size of the circular queue
    buf      unsafe.Pointer // points to an array of dataqsiz elements
    elemsize uint16
    closed   uint32
    elemtype *_type // element type
    sendx    uint   // send index
    recvx    uint   // receive index
    recvq    waitq  // list of recv waiters
    sendq    waitq  // list of send waiters
    lock mutex
}
```

Channel is a pointer to a hchan struct.

`buf` is a circular queue, `sendx` and `recvx` are the indexes of the queue.

#### Sending

- acquire the lock
- add COPY of the value to the queue
- release the lock

#### Receiving

- acquire the lock
- read the value from the queue
- release the lock

#### Buffer overflow

Sender will be paused if channel is full.

- `ch <- data` = `gopark()` -> scheduler -> change goroutine state to waiting -> free os thread

- add to `sendq` (`waitq` struct, which is a linked list of goroutines waiting for send)

```go
type waitq struct {
	first *sudog
	last  *sudog
}


type sudog struct {
	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer
    //... 
}
```

- taking data which is under the `sendx` index, reciever will check for waiting goroutines which 
are in `sendq` and will wake them up (with putting data to queue)

wake up = call of `goready()` -> scheduler -> set runnable state -> put to run queue

mutex locked only once 

Reader works in the same way with `recvq`.

`sendDirect()` will send data directly if channel is empty and some goroutines are waiting for data.
Sending direct from stack from one goroutine to another.

#### Select statement

```go
select {
case <-ch1:
    // ...
default:
    // ...
}
```

`<-ch1` = `chanrecv(hchan *chanType, ep unsafe.Pointer, block bool)` -> block = `false` because we using select 

#### Closing

- is channel initialized? (panic if not)
- lock the mutex 
- is channel closed? (panic if yes)
- set closed flag
- release all readers 
- release all writers, will panic 
- release the mutex
- unblock all goroutines

#### Sources

- [Go channels RU](https://www.youtube.com/watch?v=ZTJcaP4G4JM)
