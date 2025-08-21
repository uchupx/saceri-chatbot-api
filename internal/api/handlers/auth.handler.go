package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/internal/api/handlers/type"
	"github.com/uchupx/saceri-chatbot-api/internal/config"
	"github.com/uchupx/saceri-chatbot-api/internal/service"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
	"github.com/uchupx/saceri-chatbot-api/pkg/grpc/client"
	"github.com/uchupx/saceri-chatbot-api/pkg/grpc/proto/gen/authservice"
)

type AuthHandler struct {
	Handler
	AuthClient  *client.AuthClient
	UserService *service.UserService
}

func (h *AuthHandler) Register(c echo.Context) error {
	body := _type.RegisterRequest{}

	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	conf := config.GetConfig()
	payload := authservice.RegisterUserRequest{
		Username: body.Username,
		Password: body.Password,
		Secret:   conf.App.Secret,
	}

	respAuth, err := h.AuthClient.Register(c.Request().Context(), &payload)
	if err != nil {
		return h.responseError(c, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	body.OauthID = respAuth.Id

	user, err := h.UserService.Register(c.Request().Context(), body)
	if err != nil {
		return h.responseError(c, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	return h.responseSuccess(c, 200, user)
}

func (h *AuthHandler) Login(c echo.Context) error {

	body := _type.LoginRequest{}
	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	conf := config.GetConfig()

	payload := authservice.LoginRequest{
		Username: body.Username,
		Password: body.Password,
		Secret:   conf.App.Secret,
	}

	respAuth, err := h.AuthClient.Login(c.Request().Context(), &payload)
	if err != nil {
		return h.responseError(c, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	return h.responseSuccess(c, 200, _type.LoginResponse{
		Token: respAuth.Token,
		Exp:   3600,
	})
}
