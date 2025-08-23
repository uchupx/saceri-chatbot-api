package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
	"github.com/uchupx/saceri-chatbot-api/pkg/apilog"
)

type Handler struct {
	log *apilog.ApiLog
}

func NewHandler(log *apilog.ApiLog) *Handler {
	return &Handler{
		log: log,
	}
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Meta    any    `json:"meta,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (h Handler) responseError(c echo.Context, ctx context.Context, err *apierror.APIerror) error {
	h.log.Error(ctx, "Error occurred", err.Unwrap(), map[string]interface{}{
		"code":    err.Code(),
		"message": err.Error(),
	})

	response := Response{
		Status:  err.Code(),
		Message: err.Error(),
	}

	return c.JSON(err.Code(), response)
}

func (h Handler) responseSuccess(c echo.Context, code int, data any) error {
	return c.JSON(code, Response{
		Status:  code,
		Data:    data,
		Message: "Success",
	})
}
