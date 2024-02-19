package error

import (
	"fmt"

	"github.com/alsve/SawitProUserService/generated"
)

// StdError is base error.
type StdError struct {
	generated.StdErrResponse
	httpStatusCode int
}

// Error implements error.
func (s *StdError) Error() string {
	return fmt.Sprintf("[%s] %s", s.RequestId, s.Message)
}

// StatusCode returns http status code that represents the error.
func (s *StdError) StatusCode() int {
	return s.httpStatusCode
}

// SetRequestID sets request id.
func (s *StdError) SetRequestID(id interface{}) *StdError {
	s.RequestId = id.(string)
	return s
}
