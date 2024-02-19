package handler

import (
	"net/http"

	"github.com/alsve/SawitProUserService/generated"
	stderr "github.com/alsve/SawitProUserService/internal/error"
	stdresp "github.com/alsve/SawitProUserService/internal/stdresp"
	"github.com/alsve/SawitProUserService/repository"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Registration register a new user.
func (s *Server) Registration(ctx echo.Context) error {
	req := &generated.RegistrationJSONBody{}
	if err := ctx.Bind(req); err != nil {
		return stderr.NewBadRequest("your request is incomprehensible")
	}

	if err := ctx.Validate(req); err != nil {
		return stderr.FromValidationErrors(err.(govalidator.Errors))
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	dbin := repository.RegisterUserInput{
		Principal:        req.Principal,
		RegistrationType: req.RegistrationType,
		FullName:         req.FullName,
		PasswordMethod:   `bcrypt`,
		PasswordHash:     string(hash),
		PasswordSalt:     string(hash),
	}
	dbout, err := s.Repository.RegisterUser(ctx.Request().Context(), dbin)
	if err != nil {
		return err
	}

	return stdresp.SendJSONStdResponse[generated.RegisterResponse](
		ctx,
		http.StatusCreated,
		"user successfuly registered",
		generated.RegisterResponse{UserId: dbout.UserID},
	)
}
