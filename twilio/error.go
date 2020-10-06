package twilio

import "fmt"

// Error error
type Error struct {
	StatusCode int
	Message    string
	Detail     string
	Type       string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}
