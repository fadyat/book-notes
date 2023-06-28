package main

import (
	"context"
	"log"
	"net/http"
	"sync"
)

func serverCallWithCancel(ctx context.Context, url string, respChan chan *http.Response) {
	// Making raw request to a slow_server.go
	// Pushing the response to a channel, if the request isn't cancelled

	resp, err := serverCall(ctx, url)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	select {
	case <-ctx.Done():
		log.Printf("Request %s cancelled: %s", url, ctx.Err())
	case respChan <- resp:
		log.Printf("Pushing %s response to channel", url)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	respChan := make(chan *http.Response)
	defer close(respChan)

	baseURL := "http://localhost:8080"
	params := []string{"alice", "bob", "charlie"}

	var wg sync.WaitGroup
	for _, param := range params {
		wg.Add(1)
		url := baseURL + "?name=" + param

		go func(url string) {
			serverCallWithCancel(ctx, url, respChan)
			log.Printf("Finished request to %s", url)
			wg.Done()
		}(url)
	}

	var firstResponse *http.Response
	go func() {
		firstResponse = <-respChan
		cancel()
	}()

	wg.Wait()
	log.Printf("First response: %s", firstResponse.Status)
}
