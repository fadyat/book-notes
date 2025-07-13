package main

import (
	"context"
	"golang.org/x/time/rate"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

func (m *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range m.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *multiLimiter) Limit() rate.Limit {
	return m.limiters[0].Limit()
}

func newMultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}

	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

type apiConnection struct {
	networkLimit,
	diskLimit,
	apiLimit RateLimiter
}

func open() *apiConnection {
	return &apiConnection{
		apiLimit: newMultiLimiter(
			rate.NewLimiter(per(2, time.Second), 1),
			// with a burstiness of 10 to give the users their initial pool.
			rate.NewLimiter(per(10, time.Minute), 10),
		),
		diskLimit: newMultiLimiter(
			rate.NewLimiter(rate.Limit(1), 1),
		),
		networkLimit: newMultiLimiter(
			rate.NewLimiter(per(3, time.Second), 3),
		),
	}
}

func (a *apiConnection) resolveAddress(ctx context.Context) error {
	limiter := newMultiLimiter(a.apiLimit, a.networkLimit)

	if err := limiter.Wait(ctx); err != nil {
		return err
	}

	// work here...

	return nil
}

func (a *apiConnection) readFile(ctx context.Context) error {
	limiter := newMultiLimiter(a.apiLimit, a.diskLimit)

	if err := limiter.Wait(ctx); err != nil {
		return err
	}

	// work here...

	return nil
}

func per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}

func main() {
	defer log.Printf("done")

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConn := open()

	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			if err := apiConn.readFile(context.Background()); err != nil {
				log.Printf("can't readfile: %v", err)
			}

			log.Printf("read file")
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			if err := apiConn.resolveAddress(context.Background()); err != nil {
				log.Printf("can't resolve address: %v", err)
			}

			log.Printf("resolve address")
		}()
	}

	wg.Wait()
}
