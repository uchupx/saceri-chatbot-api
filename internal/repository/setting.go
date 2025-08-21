package repository

import "github.com/uchupx/saceri-chatbot-api/internal/models"

type SettingRepoInterface interface {
	Create(data models.SettingModel) (*models.SettingModel, error)
	Update(data models.SettingModel) (*models.SettingModel, error)
	GetByKey(key models.SettingKey) (*models.SettingModel, error)
	GetAllSettings() ([]models.SettingModel, error)
}
