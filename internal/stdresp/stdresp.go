package stdresp

import (
	"github.com/alsve/SawitProUserService/generated"
	"github.com/alsve/SawitProUserService/internal/middleware"
	"github.com/labstack/echo/v4"
)

// SendJSONStdResponse send data via ctx.JSON.
func SendJSONStdResponse[Data any](ctx echo.Context, httpCode int, message string, data Data) error {
	reqId := ""
	if rid, ok := ctx.Get(middleware.CKeyReqID).(string); ok {
		reqId = rid
	}

	return ctx.JSON(httpCode, &generated.StdResponse{
		Status:    "success",
		RequestId: reqId,
		Message:   message,
		Data:      data,
	})
}
