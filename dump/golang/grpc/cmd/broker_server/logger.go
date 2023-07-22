package main

import (
	"log"
	"os"
)

func initLogger() *log.Logger {
	return log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}
