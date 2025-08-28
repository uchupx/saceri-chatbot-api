package repository

import (
	"context"

	"github.com/uchupx/saceri-chatbot-api/internal/models"
)

type EventRepoInterface interface {
	Create(ctx context.Context, event models.EventModel) (*models.EventModel, error)
	Update(ctx context.Context, event models.EventModel) (*models.EventModel, error)
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (*models.EventModel, error)
	GetAllEvents(ctx context.Context, limit, offset int) ([]models.EventModel, error)
}
