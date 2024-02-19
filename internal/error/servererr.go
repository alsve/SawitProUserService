package error

import (
	"fmt"
	"net/http"

	"github.com/alsve/SawitProUserService/generated"
)

const (
	ErrInsertUserSQLStmt          = "000001"
	ErrInsertUserScan             = "000002"
	ErrInsertUserCommit           = "000003"
	ErrGetPasswordByUserIDSQLStmt = "000004"
	ErrGetPasswordByUserIDScan    = "000005"
	ErrAuthSignString             = "000006"
)

type ServerError struct {
	StdError
	rootErr   error
	errorCode string
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", s.RequestId, s.errorCode, s.rootErr.Error())
}

// NewServerError creates ServerError by error code and the root error that cause the problem.
func NewServerError(errorCode string, err error) *ServerError {
	message := fmt.Sprintf("%s: Internal Server Error", errorCode)
	serr := &ServerError{
		errorCode: errorCode,
		rootErr:   err,
		StdError: StdError{
			httpStatusCode: http.StatusInternalServerError,
			StdErrResponse: generated.StdErrResponse{
				Status:  StatusFailed,
				Message: message,
				Code:    CodeInternalServerError,
				Data:    nil,
			},
		},
	}
	return serr
}
