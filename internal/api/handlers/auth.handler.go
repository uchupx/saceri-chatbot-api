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

	ctx := h.log.CreateTrace(c.Request().Context())
	body := _type.RegisterRequest{}

	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, body)

	conf := config.GetConfig()
	payload := authservice.RegisterUserRequest{
		Username: body.Username,
		Password: body.Password,
		Secret:   conf.App.Secret,
	}

	respAuth, err := h.AuthClient.Register(ctx, &payload)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	body.OauthID = respAuth.Id

	user, err := h.UserService.Register(ctx, body)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	return h.responseSuccess(c, 200, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	body := _type.LoginRequest{}

	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, body)
	conf := config.GetConfig()

	payload := authservice.LoginRequest{
		Username: body.Username,
		Password: body.Password,
		Secret:   conf.App.Secret,
	}

	respAuth, err := h.AuthClient.Login(c.Request().Context(), &payload)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	return h.responseSuccess(c, 200, _type.LoginResponse{
		Token: respAuth.Token,
		Exp:   3600,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {

	token := c.Get("token").(string)
	ctx := h.log.CreateTrace(c.Request().Context())

	payload := authservice.LogoutRequest{
		Token: token,
	}

	_, err := h.AuthClient.Logout(ctx, &payload)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrInternalServerError.Code, err))
	}

	return h.responseSuccess(c, 200, "success logout")
}
