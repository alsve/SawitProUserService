package handler

import (
	"net/http"
	"time"

	"github.com/alsve/SawitProUserService/generated"
	stderr "github.com/alsve/SawitProUserService/internal/error"
	stdresp "github.com/alsve/SawitProUserService/internal/stdresp"
	"github.com/alsve/SawitProUserService/repository"
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Authentication authenticate user and returns jwt token and expire time.
func (s *Server) Authentication(ctx echo.Context) error {
	req := &generated.AuthenticationJSONBody{}
	if err := ctx.Bind(req); err != nil {
		return stderr.NewBadRequest("your request is incomprehensible")
	}

	if err := ctx.Validate(req); err != nil {
		return stderr.FromValidationErrors(err.(govalidator.Errors))
	}

	dbin := repository.GetAuthDataInput{
		Principal:        req.Principal,
		RegistrationType: req.RegistrationType,
	}
	dbout, err := s.Repository.GetAuthData(ctx.Request().Context(), dbin)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbout.PasswordHash), []byte(req.Password))
	if err != nil {
		return stderr.NewBadRequest("wrong user or password")
	}

	claim := jwt.MapClaims{
		"sub": dbout.UserID,
		"aud": "sawit_pro_user_service",
		"iss": "sawit_pro_user_service",
		"exp": time.Now().Add(time.Duration(s.defaultJWTExpirationInSecond) * time.Second).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	accessToken, err := token.SignedString(s.jwtRSAPrivateKey)
	if err != nil {
		return stderr.NewServerError(stderr.ErrAuthSignString, err)
	}

	return stdresp.SendJSONStdResponse[generated.LoginResponse](
		ctx,
		http.StatusOK,
		"successfully authenticating user",
		generated.LoginResponse{
			TokenType:   "Bearer",
			AccessToken: accessToken,
			ExpiresIn:   s.defaultJWTExpirationInSecond,
		},
	)
}
