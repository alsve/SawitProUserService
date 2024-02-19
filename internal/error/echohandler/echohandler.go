package echohandler

import (
	"errors"

	stderr "github.com/alsve/SawitProUserService/internal/error"
	"github.com/alsve/SawitProUserService/internal/middleware"

	"github.com/labstack/echo/v4"
)

// NewEchoErrorHandler creates a new instance of EchoHandler.
func NewEchoErrorHandler() *EchoErrorHandler {
	return &EchoErrorHandler{}
}

// EchoErrorHandler is used to handles error that implements interface of
// echo.HTTPErrorHandler.
type EchoErrorHandler struct{}

func (e *EchoErrorHandler) ErrorHandler(err error, ctx echo.Context) {
	var cerr *stderr.ClientError
	if errors.As(err, &cerr) {
		cerr.SetRequestID(ctx.Get(middleware.CKeyReqID))
		ctx.JSON(cerr.StatusCode(), cerr.StdErrResponse)
		return
	}

	var serr *stderr.ServerError
	if errors.As(err, &serr) {
		serr.SetRequestID(ctx.Get(middleware.CKeyReqID))
		ctx.Logger().Error(serr)
		ctx.JSON(serr.StatusCode(), serr.StdErrResponse)
	}
}
