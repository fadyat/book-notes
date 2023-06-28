package main

import (
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Sleeping for 5 seconds to simulate a workload
	time.Sleep(5 * time.Second)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello, " + r.URL.Query().Get("name")))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Started request")
		next.ServeHTTP(w, r)
		log.Println("Finished handling request")
	})
}

func main() {
	server := &http.Server{
		Addr:        ":8080",
		Handler:     loggerMiddleware(http.HandlerFunc(handler)),
		ReadTimeout: 2 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
