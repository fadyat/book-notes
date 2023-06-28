package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	log.Printf("Storing `%s`: `%s` in context", "key", "value")
	valueCtx := context.WithValue(ctx, "key", "value")

	value := valueCtx.Value("key")
	log.Printf("Key `%s` has value `%s`", "key", value)
}
