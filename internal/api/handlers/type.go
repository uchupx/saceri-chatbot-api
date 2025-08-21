package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
)

type Handler struct {
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Meta    any    `json:"meta,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (h Handler) responseError(c echo.Context, err *apierror.APIerror) error {
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
