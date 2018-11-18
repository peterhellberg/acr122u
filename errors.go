package acr122u

// Error is tye error type returned by the acr122u package
type Error struct {
	Message string
}

// Error returns the error message
func (e Error) Error() string {
	return e.Message
}

var (
	// ErrOperationFailed is returned when the response code is 0x63 0x00
	ErrOperationFailed = Error{Message: "operation failed"}
)
