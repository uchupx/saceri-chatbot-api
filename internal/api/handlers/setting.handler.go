package handlers

import (
	"github.com/labstack/echo/v4"
	_type "github.com/uchupx/saceri-chatbot-api/internal/api/handlers/type"
	"github.com/uchupx/saceri-chatbot-api/internal/service"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
)

type SettingHandler struct {
	Handler
	SettingService *service.SettingService
}

func (h *SettingHandler) GetSetting(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	settings, err := h.SettingService.GetAll(ctx)
	if err != nil {
		return h.responseError(c, ctx, err)
	}

	return h.responseSuccess(c, 200, settings)
}

func (h *SettingHandler) UpdateSetting(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	body := _type.SettingUpdateRequest{}

	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, body)

	setting, apierr := h.SettingService.Update(ctx, body.Key, body.Value)
	if apierr != nil {
		return h.responseError(c, ctx, apierr)
	}

	return h.responseSuccess(c, 200, setting)
}
