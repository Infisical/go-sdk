package errors

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/infisical/go-sdk/packages/util"
)

// APIError represents an error response from the API
type APIError struct {
	Operation    string  `json:"operation"`
	Method       string  `json:"method"`
	URL          string  `json:"url"`
	StatusCode   int     `json:"status_code"`
	ErrorMessage *string `json:"response,omitempty"`
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf(
		"APIError: %s unsuccessful response [%v %v] [status-code=%v]",
		e.Operation,
		e.Method,
		e.URL,
		e.StatusCode,
	)

	if e.ErrorMessage != nil {
		return fmt.Sprintf("%s Error: %s", msg, *e.ErrorMessage)

	}

	return msg
}

func NewAPIError(operation string, res *resty.Response) error {
	return &APIError{
		Operation:  operation,
		Method:     res.Request.Method,
		URL:        res.Request.URL,
		StatusCode: res.StatusCode(),
	}
}

func NewAPIErrorWithResponse(operation string, res *resty.Response) error {
	errorMessage := util.TryParseErrorBody(res)

	return &APIError{
		Operation:    operation,
		Method:       res.Request.Method,
		URL:          res.Request.URL,
		StatusCode:   res.StatusCode(),
		ErrorMessage: &errorMessage,
	}
}

func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}
