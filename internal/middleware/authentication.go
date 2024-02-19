package middleware

import (
	"crypto/rsa"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// Contains data related to authentication.
type AuthenticationData struct {
	UserID string
}

// NewAuthenticationMiddleware creates a new instance of AuthenticationMiddleware.
func NewAuthenticationMiddleware(JWTRSAPrivateKey *rsa.PrivateKey) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{jwtRSAPrivateKey: JWTRSAPrivateKey}
}

// AuthenticationMiddleware checks for integrity of access_token (JWT) from authorisation header, checks for expiration and
// user data existance and transform it to authentication data.
type AuthenticationMiddleware struct {
	jwtRSAPrivateKey *rsa.PrivateKey
}

func (a *AuthenticationMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return a.authenticate(c, next)
	}
}

func (a *AuthenticationMiddleware) authenticate(ctx echo.Context, next echo.HandlerFunc) (err error) {
	separator := " "
	authStr := ctx.Request().Header.Get(echo.HeaderAuthorization)
	if !strings.Contains(authStr, separator) {
		return next(ctx)
	}

	tokenStr := strings.Split(authStr, separator)[1]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return a.jwtRSAPrivateKey.Public(), nil
	})
	if err != nil {
		return next(ctx)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return next(ctx)
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return next(ctx)
	}

	ctx.Set(CKeyAuthData, AuthenticationData{
		UserID: subject,
	})

	return next(ctx)
}
