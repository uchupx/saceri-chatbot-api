package handlers

import (
	"github.com/labstack/echo/v4"
	_type "github.com/uchupx/saceri-chatbot-api/internal/api/handlers/type"
	"github.com/uchupx/saceri-chatbot-api/internal/service"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
)

type EventHandler struct {
	Handler
	EventService *service.EventService
}

func (h *EventHandler) CreateEvent(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	body := _type.EventCreateUpdateRequest{}

	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, body)

	event, apierr := h.EventService.Create(ctx, body)
	if apierr != nil {
		return h.responseError(c, ctx, apierr)
	}

	return h.responseSuccess(c, 201, event)
}

func (h *EventHandler) GetEvent(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	id := c.Param("id")

	event, err := h.EventService.GetById(ctx, id)
	if err != nil {
		return h.responseError(c, ctx, err)
	}

	return h.responseSuccess(c, 200, event)
}

func (h *EventHandler) UpdateEvent(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	id := c.Param("id")
	body := _type.EventCreateUpdateRequest{}

	err := c.Bind(&body)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, body)

	event, apierr := h.EventService.Update(ctx, id, body)
	if apierr != nil {
		return h.responseError(c, ctx, apierr)
	}

	return h.responseSuccess(c, 200, event)
}

func (h *EventHandler) DeleteEvent(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())
	id := c.Param("id")

	err := h.EventService.Delete(ctx, id)
	if err != nil {
		return h.responseError(c, ctx, err)
	}

	return h.responseSuccess(c, 204, nil)
}

func (h *EventHandler) GetEvents(c echo.Context) error {
	ctx := h.log.CreateTrace(c.Request().Context())

	query := _type.GetQuery{}
	err := c.Bind(&query)
	if err != nil {
		return h.responseError(c, ctx, apierror.NewAPIError(echo.ErrBadRequest.Code, err))
	}

	ctx = h.log.AttachBody(ctx, query)

	events, apierr := h.EventService.GetAll(ctx, query)
	if apierr != nil {
		return h.responseError(c, ctx, apierr)
	}

	return h.responseSuccess(c, 200, events)
}
