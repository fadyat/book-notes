package logger

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"log"
)

func ToInterceptorLogger(l *log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, args ...interface{}) {
		switch lvl {
		case logging.LevelDebug:
			msg = "debug: " + msg
		case logging.LevelInfo:
			msg = "info: " + msg
		case logging.LevelWarn:
			msg = "warn: " + msg
		case logging.LevelError:
			msg = "error: " + msg
		}

		l.Printf(msg, args...)
	})
}
