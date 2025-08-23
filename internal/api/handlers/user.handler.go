package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/internal/service"
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

	users, err := h.UserService.GetUsers(ctx)
	if err != nil {
		return h.responseError(c, ctx, err)
	}
	return h.responseSuccess(c, 200, users)
}
