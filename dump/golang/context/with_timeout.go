package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func serverCallWithChannel(ctx context.Context, url string, respChan chan<- *http.Response) {
	res, err := serverCall(ctx, url)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	respChan <- res
}

func serverCallMoreTimeout(ctx context.Context, url string, respChan chan<- *http.Response) {
	// We can create a new context from the parent context
	// Can't overwrite the parent context value, but can add new values
	// Timeout still be 2 seconds

	moreCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	serverCallWithChannel(moreCtx, url, respChan)
}

func main() {
	// Imagine our application is accessing a slow external service
	//
	// We're not sure how long we should wait for a response
	// so we don't block the app -> setting a timeout for the request via context

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	respChan := make(chan *http.Response)
	defer close(respChan)

	//go serverCall(timeoutCtx, respChan)
	go serverCallMoreTimeout(ctx, "http://localhost:8080", respChan)

	select {
	case <-ctx.Done():
		log.Printf("Request timed out: %s", ctx.Err())
	case res := <-respChan:
		log.Printf("Response status: %s", res.Status)
	}
}
