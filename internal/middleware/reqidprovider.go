package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NewRequestIDProvider creates a new instance of RequestIDProvider.
func NewRequestIDProvider() *RequestIDProvider {
	return &RequestIDProvider{}
}

// RequestIDProvider provide request id for the current request.
type RequestIDProvider struct{}

func (r *RequestIDProvider) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return r.middleware(ctx, next)
	}
}

func (r *RequestIDProvider) middleware(ctx echo.Context, next echo.HandlerFunc) error {
	ctx.Set(CKeyReqID, uuid.New().String())
	return next(ctx)
}
