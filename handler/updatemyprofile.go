package handler

import (
	"net/http"

	"github.com/alsve/SawitProUserService/generated"
	stderr "github.com/alsve/SawitProUserService/internal/error"
	"github.com/alsve/SawitProUserService/internal/middleware"
	"github.com/alsve/SawitProUserService/internal/stdresp"
	"github.com/alsve/SawitProUserService/repository"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

// UpdateMyProfile updates user profile.
// User should get authenticated first.
func (s *Server) UpdateMyProfile(ctx echo.Context) error {
	authData, ok := ctx.Get(middleware.CKeyAuthData).(middleware.AuthenticationData)
	if !ok {
		return stderr.NewForbidden("you are not authorized to do that")
	}

	req := &generated.UpdateMyProfileJSONBody{}
	if err := ctx.Bind(req); err != nil {
		return stderr.NewBadRequest("your request is incomprehensible")
	}

	if err := ctx.Validate(req); err != nil {
		return stderr.FromValidationErrors(err.(govalidator.Errors))
	}

	dbin := repository.UpdateUserInput{
		UserID:      authData.UserID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}
	dbout, err := s.Repository.UpdateUser(ctx.Request().Context(), dbin)
	if err != nil {
		return err
	}

	return stdresp.SendJSONStdResponse[generated.UserInfoResponse](ctx, http.StatusOK, "update user successful", generated.UserInfoResponse{
		FullName:  dbout.FullName,
		Principal: dbout.PhoneNumber,
	})
}
