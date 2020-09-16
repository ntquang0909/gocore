package sendgrid

// ErrorResponse errors
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error error
type Error struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Help    string `json:"help"`
}

// Error implement error func
func (e Error) Error() string {
	return e.Message
}
