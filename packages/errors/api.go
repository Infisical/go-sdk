package errors

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/levidurfee/go-sdk/packages/util"
)

// APIError represents an error response from the API
type APIError struct {
	Operation    string `json:"operation"`
	Method       string `json:"method"`
	URL          string `json:"url"`
	StatusCode   int    `json:"statusCode"`
	ErrorMessage string `json:"message,omitempty"`
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf(
		"APIError: %s unsuccessful response [%v %v] [status-code=%v]",
		e.Operation,
		e.Method,
		e.URL,
		e.StatusCode,
	)

	if e.ErrorMessage != "" {
		return fmt.Sprintf("%s [message=\"%s\"]", msg, e.ErrorMessage)

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
		ErrorMessage: errorMessage,
	}
}
