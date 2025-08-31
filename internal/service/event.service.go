package service

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	_type "github.com/uchupx/saceri-chatbot-api/internal/api/handlers/type"
	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
	"github.com/uchupx/saceri-chatbot-api/pkg/helper"
)

type EventService struct {
	Repo           repository.EventRepoInterface
	ChatbotService *ChatbotService
}

func (e *EventService) Create(ctx context.Context, evt _type.EventCreateUpdateRequest) (*models.EventModel, *apierror.APIerror) {
	event := models.EventModel{
		Name:           evt.Title,
		Description:    evt.Description,
		Place:          evt.Place,
		DonationTarget: evt.DonationTarget,
		StartEvent:     helper.StringToTime(evt.StartAt),
		EndEvent:       helper.StringToTime(evt.EndAt),
	}

	res, err := e.Repo.Create(ctx, event)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	go func() {
		ctx := context.Background()
		err := e.ChatbotService.UpdatePromptContext(ctx)
		if err != nil {
			fmt.Printf("Error updating prompt context: %v\n", err.Error())
			return
		}
	}()

	return res, nil

}

func (e *EventService) Update(ctx context.Context, id string, evt _type.EventCreateUpdateRequest) (*models.EventModel, *apierror.APIerror) {
	event, err := e.Repo.GetById(ctx, id)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if event == nil {
		return nil, apierror.NewAPIError(echo.ErrNotFound.Code, apierror.ERR_NOT_FOUND)
	}

	event.Name = evt.Title
	event.Description = evt.Description
	event.Place = evt.Place
	event.DonationTarget = evt.DonationTarget
	event.StartEvent = helper.StringToTime(evt.StartAt)
	event.EndEvent = helper.StringToTime(evt.EndAt)
	_, err = e.Repo.Update(ctx, *event)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	go func() {
		ctx := context.Background()
		err := e.ChatbotService.UpdatePromptContext(ctx)
		if err != nil {
			fmt.Printf("Error updating prompt context: %v\n", err.Error())
			return
		}
	}()

	return event, nil
}

func (e *EventService) Delete(ctx context.Context, id string) *apierror.APIerror {
	err := e.Repo.Delete(ctx, id)
	if err != nil {
		return apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return nil
}

func (e *EventService) GetById(ctx context.Context, id string) (*models.EventModel, *apierror.APIerror) {
	event, err := e.Repo.GetById(ctx, id)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if event == nil {
		return nil, apierror.NewAPIError(echo.ErrNotFound.Code, apierror.ERR_NOT_FOUND)
	}

	return event, nil
}

func (e *EventService) GetAll(ctx context.Context, params _type.GetQuery) ([]models.EventModel, *apierror.APIerror) {
	events, err := e.Repo.GetAllEvents(ctx, params.Keyword, params.Limit(), params.Offset())
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return events, nil
}

func (e *EventService) GetByDates(ctx context.Context, datesString []string) ([]models.EventModel, *apierror.APIerror) {
	var dates []time.Time

	for _, dateStr := range datesString {
		date := helper.StringToTime(dateStr)
		dates = append(dates, date)
	}

	events, err := e.Repo.GetEventsByDates(ctx, dates)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return events, nil
}
