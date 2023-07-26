package main

import (
	"log"
	"strings"
)

func mayBeDie(err error, msg ...string) {
	if err != nil {
		log.Fatalf("%s: %s", strings.Join(msg, ""), err)
	}
}
