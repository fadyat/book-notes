### What is context

From docs:

```go
// Package context defines the Context type, which carries deadlines,
// cancellation signals, and other request-scoped values across API boundaries
// and between processes.
```

#### Some context types:

```go
background = new(emptyCtx)

// Background returns a non-nil, empty Context. It is never canceled, has no
// values, and has no deadline. It is typically used by the main function,
// initialization, and tests, and as the top-level Context for incoming
// requests.
func Background() Context {
    return background
}
```

```go
todo = new(emptyCtx)

// TODO returns a non-nil, empty Context. Code should use context.TODO when
// it's unclear which Context to use or it is not yet available (because the
// surrounding function has not yet been extended to accept a Context
// parameter).
func TODO() Context {
    return todo
}
```

```go
// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct {
    Context
    key, val any
}

// WithValue returns a copy of parent in which the value associated with key is
// val.
//
// Use context Values only for request-scoped data that transits processes and
// APIs, not for passing optional parameters to functions.
//
// The provided key must be comparable and should not be of type
// string or any other built-in type to avoid collisions between
// packages using context. Users of WithValue should define their own
// types for keys. To avoid allocating when assigning to an
// interface{}, context keys often have concrete type
// struct{}. Alternatively, exported context key variables' static
// type should be a pointer or interface.
func WithValue(parent Context, key, val any) Context {
    // ...

    return &valueCtx{parent, key, val}
}
```

```go
// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete:
//
//  func slowOperationWithTimeout(ctx context.Context) (Result, error) {
//      ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
//      defer cancel()  // releases resources if slowOperation completes before timeout elapses
//      return slowOperation(ctx)
//  }
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
    return WithDeadline(parent, time.Now().Add(timeout))
}
```

```go
// A timerCtx carries a timer and a deadline. It embeds a cancelCtx to
// implement Done and Err. It implements cancel by stopping its timer then
// delegating to cancelCtx.cancel.
type timerCtx struct {
    *cancelCtx
    timer *time.Timer // Under cancelCtx.mu.

    deadline time.Time
}

// WithDeadline returns a copy of the parent context with the deadline adjusted
// to be no later than d. If the parent's deadline is already earlier than d,
// WithDeadline(parent, d) is semantically equivalent to parent. The returned
// context's Done channel is closed when the deadline expires, when the returned
// cancel function is called, or when the parent context's Done channel is
// closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
    // ...
    c := &timerCtx{
        cancelCtx: newCancelCtx(parent),
        deadline:  d,
    }

    // ...
    return c, func() { c.cancel(true, Canceled, nil) }
}
```

```go
// WithCancel returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called
// or when the parent context's Done channel is closed, whichever happens first.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
    c := withCancel(parent)
    return c, func() { c.cancel(true, Canceled, nil) }
}
```

### Best practices:

- Pass `context.Context` as the first argument to every function that needs it.
- Pass only to a functions, not to a struct fields.
- Pass only `context.Context` without cancel function.
- Use `context.WithValue` only in extreme cases.
- Use `context.Background` only as a top-level context.
- Use `context.TODO` only when you don't know which context to use.
- Never pass a `nil` context, use `context.TODO` instead.

Documentation: [https://pkg.go.dev/context](https://pkg.go.dev/context)
