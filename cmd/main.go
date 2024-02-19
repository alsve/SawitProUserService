package main

import (
	"crypto/rsa"
	"os"
	"regexp"
	"strconv"

	"github.com/alsve/SawitProUserService/generated"
	"github.com/alsve/SawitProUserService/handler"
	"github.com/alsve/SawitProUserService/internal/error/echohandler"
	"github.com/alsve/SawitProUserService/internal/middleware"
	"github.com/alsve/SawitProUserService/repository"
	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
)

// initiate govalidator extension if there is any.
func newCustomValidator() *customValidator {
	phoneNumberRegexp := regexp.MustCompile(`^\+628[0-9]+$`)
	hasLowerCase := regexp.MustCompile(`[a-z]+`)
	hasUpperCase := regexp.MustCompile(`[A-Z]+`)
	hasNumber := regexp.MustCompile(`[0-9]+`)
	hasSpecialCharacter := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	govalidator.TagMap["phone-number-pattern"] = func(str string) bool {
		return phoneNumberRegexp.MatchString(str)
	}
	govalidator.TagMap["has-lower-case"] = func(str string) bool {
		return hasLowerCase.MatchString(str)
	}
	govalidator.TagMap["has-upper-case"] = func(str string) bool {
		return hasUpperCase.MatchString(str)
	}
	govalidator.TagMap["has-number"] = func(str string) bool {
		return hasNumber.MatchString(str)
	}
	govalidator.TagMap["has-special-character"] = func(str string) bool {
		return hasSpecialCharacter.MatchString(str)
	}
	return &customValidator{}
}

// customValidator is an adapter for go-playground validator to echo validator.
type customValidator struct{}

// Validate implements echo.Validator.
func (c *customValidator) Validate(v interface{}) error {
	_, err := govalidator.ValidateStruct(v)
	return err
}

func main() {
	e := echo.New()

	// build echo dependencies.
	eh := echohandler.NewEchoErrorHandler()
	e.HTTPErrorHandler = eh.ErrorHandler
	e.Validator = newCustomValidator()

	// build echo middlewares.
	privateKeyByte, err := os.ReadFile("jwt_private.pem")
	if err != nil {
		panic(err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyByte)
	if err != nil {
		panic(err)
	}

	reqidprovider := middleware.NewRequestIDProvider()
	optionalJWTAuth := middleware.NewAuthenticationMiddleware(privateKey)
	e.Use(optionalJWTAuth.Middleware)
	e.Use(reqidprovider.Middleware)

	var server generated.ServerInterface = newServer(privateKey)

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer(jwtPrivateKey *rsa.PrivateKey) *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	jwtExpInSec, _ := strconv.Atoi(os.Getenv("DEFAULT_JWT_EXP_IN_SEC"))

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository:                   repo,
		DefaultJWTExpirationInSecond: jwtExpInSec,
		JWTRSAPrivateKey:             jwtPrivateKey,
	}
	return handler.NewServer(opts)
}
