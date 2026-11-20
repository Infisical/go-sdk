package errors

import (
	"fmt"
)

// RequestError represents an error that occurred during an API request
type RequestError struct {
	Operation string `json:"operation"`
	error     error  `json:"-"`
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%s: unable to complete api request [err=%s]", e.Operation, e.error)
}

// Unwrap returns the error the request failed with, so callers can inspect it with errors.Is and errors.As. A
// *net.OpError whose Op is "dial", for example, means the connection never opened and the server never received the
// request.
func (e *RequestError) Unwrap() error {
	return e.error
}

func NewRequestError(operation string, err error) error {
	return &RequestError{
		Operation: operation,
		error:     err,
	}
}
