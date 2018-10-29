package acr122u

// Logger interface that is implemented by *log.Logger
type Logger interface {
	Printf(format string, v ...interface{})
}
