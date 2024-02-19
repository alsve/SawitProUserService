package handler

import (
	"net/http"

	"github.com/alsve/SawitProUserService/generated"
	stderr "github.com/alsve/SawitProUserService/internal/error"
	"github.com/alsve/SawitProUserService/internal/middleware"
	"github.com/alsve/SawitProUserService/internal/stdresp"
	"github.com/alsve/SawitProUserService/repository"
	"github.com/labstack/echo/v4"
)

// GetMyProfile retrieve user profile.
// User should get authenticated first.
func (s *Server) GetMyProfile(ctx echo.Context) error {
	authData, ok := ctx.Get(middleware.CKeyAuthData).(middleware.AuthenticationData)
	if !ok {
		return stderr.NewForbidden("you are not authorized to do that")
	}

	dbin := repository.GetUserInput{
		UserID: authData.UserID,
	}
	dbout, err := s.Repository.GetUser(ctx.Request().Context(), dbin)
	if err != nil {
		return nil
	}

	stdresp.SendJSONStdResponse[generated.UserInfoResponse](
		ctx, http.StatusOK,
		"success retreiving user info",
		generated.UserInfoResponse{
			FullName:  dbout.FullName,
			Principal: dbout.PhoneNumber,
		},
	)
	return nil
}
