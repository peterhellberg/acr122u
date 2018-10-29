package acr122u

import (
	"log"
	"os"
)

// Logger interface that is implemented by *log.Logger
type Logger interface {
	Printf(format string, v ...interface{})
}

// StdoutLogger logs to Stdout with no prefix or flags
func StdoutLogger() *log.Logger {
	return log.New(os.Stdout, "", 0)
}
