package error

import (
	"fmt"
	"net/http"

	"github.com/alsve/SawitProUserService/generated"
	"github.com/asaskevich/govalidator"
)

const (
	StatusFailed = "failed"
)

const (
	CodeInternalServerError = "internal_server_error"
	CodeBadRequest          = "bad_request"
	CodeForbidden           = "forbidden"
	CodeConflict            = "conflict"
)

// ClientErrorData is a standard client error data.
type ClientErrorData struct {
	Detail []ClientErrorDataDetailItem `json:"detail,omitempty"`
}

type ClientErrorDataDetailItem struct {
	Field   *string  `json:"field,omitempty"`
	Message string   `json:"message"`
	Path    []string `json:"path,omitempty"`
}

type ClientError struct {
	StdError
}

// FromValdiationErrors construct ClientError from ValidationErrors.
func FromValidationErrors(errs govalidator.Errors) *ClientError {
	ced := &ClientErrorData{Detail: []ClientErrorDataDetailItem{}}
	for _, e := range errs {
		err := e.(govalidator.Error)
		field := err.Name
		path := err.Path
		ced.Detail = append(ced.Detail, ClientErrorDataDetailItem{
			Field:   &field,
			Message: err.Error(),
			Path:    path,
		})
	}
	return NewGeneralClientError(http.StatusBadRequest, CodeBadRequest, ced)
}

// NewBadRequest creates a ClientError for a bad request.
func NewBadRequest(detail string) *ClientError {
	return NewGeneralClientError(
		http.StatusBadRequest,
		CodeBadRequest,
		&ClientErrorData{
			Detail: []ClientErrorDataDetailItem{{
				Message: detail,
			}},
		},
	)
}

// NewForbidden creates a ClientError for a forbidden request.
func NewForbidden(detail string) *ClientError {
	return NewGeneralClientError(
		http.StatusForbidden,
		CodeForbidden,
		&ClientErrorData{
			Detail: []ClientErrorDataDetailItem{{
				Message: detail,
			}},
		},
	)
}

// NewBadRequest creates a ClientError for conflict request.
func NewConflictRequest(detail string) *ClientError {
	return NewGeneralClientError(
		http.StatusConflict,
		CodeConflict,
		&ClientErrorData{
			Detail: []ClientErrorDataDetailItem{
				{Message: detail},
			},
		},
	)
}

// NewGeneralClientError implements general client error.
func NewGeneralClientError(httpStatusCode int, code string, data interface{}) *ClientError {
	message := fmt.Sprintf("%d: %s", httpStatusCode, code)
	return &ClientError{
		StdError: StdError{
			httpStatusCode: httpStatusCode,
			StdErrResponse: generated.StdErrResponse{
				Status:  StatusFailed,
				Message: message,
				Code:    code,
				Data:    data,
			},
		},
	}
}
