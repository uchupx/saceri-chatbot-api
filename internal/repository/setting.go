package repository

import (
	"context"

	"github.com/uchupx/saceri-chatbot-api/internal/models"
)

type SettingRepoInterface interface {
	Create(ctx context.Context, data models.SettingModel) (*models.SettingModel, error)
	Update(ctx context.Context, data models.SettingModel) (*models.SettingModel, error)
	GetByKey(ctx context.Context, key models.SettingKey) (*models.SettingModel, error)
	GetAllSettings(ctx context.Context) ([]models.SettingModel, error)
}
