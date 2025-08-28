package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"
	_type "github.com/uchupx/saceri-chatbot-api/internal/api/handlers/type"
	"github.com/uchupx/saceri-chatbot-api/internal/service"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
	"github.com/uchupx/saceri-chatbot-api/pkg/grpc/proto/gen/authservice"
)

type UserHandler struct {
	Handler
	UserService *service.UserService
}

func (h *UserHandler) GetUser(c echo.Context) error {
	return h.responseSuccess(c, 200, "User retrieved successfully")
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())

	query := _type.UserGetQuery{}
	if err := c.Bind(&query); err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, query)

	users, err := h.UserService.GetUsers(ctx, query)
	if err != nil {
		return h.responseError(c, ctx, err)
	}
	return h.responseSuccess(c, 200, users)
}

func (h *UserHandler) GetMe(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())

	userCtx, ok := c.Get("user").(*authservice.GetUserResponse)
	if !ok {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrInternalServerError.Code, fmt.Errorf("failed to get user from context")))
	}

	user, err := h.UserService.GetUserByOauthID(ctx, userCtx.Id)
	if err != nil {
		return h.responseError(c, ctx, err)
	}

	return h.responseSuccess(c, 200, user)
}

func (h *UserHandler) UpdateMe(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())

	body := _type.UserUpdateRequest{}
	if err := c.Bind(&body); err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, body)

	userCtx, ok := c.Get("user").(*authservice.GetUserResponse)
	if !ok {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrInternalServerError.Code, fmt.Errorf("failed to get user from context")))
	}

	body.SetID(userCtx.Id)

	user, err := h.UserService.UpdateUser(ctx, body)
	if err != nil {
		return h.responseError(c, ctx, err)
	}

	return h.responseSuccess(c, 200, user)
}
