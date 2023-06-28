package main

import (
	"context"
	"log"
	"net/http"
)

func serverCall(ctx context.Context, url string) (*http.Response, error) {
	// Making a request to a slow_server.go
	// Request will timeout because of the context

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	log.Printf("Making request to %s", req.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
