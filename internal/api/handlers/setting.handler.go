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
	settings, err := h.SettingService.GetAll(c.Request().Context())
	if err != nil {
		return h.responseError(c, err)
	}

	return h.responseSuccess(c, 200, settings)
}

func (h *SettingHandler) UpdateSetting(c echo.Context) error {
	body := _type.SettingUpdateRequest{}
	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	setting, apierr := h.SettingService.Update(c.Request().Context(), body.Key, body.Value)
	if apierr != nil {
		return h.responseError(c, apierr)
	}

	return h.responseSuccess(c, 200, setting)
}
