package service

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
	"github.com/uchupx/saceri-chatbot-api/pkg/apierror"
)

type SettingService struct {
	SettingRepo repository.SettingRepoInterface
}

func (s *SettingService) Update(ctx context.Context, key models.SettingKey, value string) (*models.SettingModel, *apierror.APIerror) {
	setting, err := s.SettingRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if setting == nil {
		return nil, apierror.NewAPIError(echo.ErrNotFound.Code, nil)
	}

	setting.Value = value
	_, err = s.SettingRepo.Update(ctx, *setting)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return setting, nil
}

func (s *SettingService) GetByKey(ctx context.Context, key models.SettingKey) (*models.SettingModel, *apierror.APIerror) {
	setting, err := s.SettingRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	if setting == nil {
		return nil, apierror.NewAPIError(echo.ErrNotFound.Code, nil)
	}

	return setting, nil
}

func (s *SettingService) GetAll(ctx context.Context) ([]models.SettingModel, *apierror.APIerror) {
	settings, err := s.SettingRepo.GetAllSettings(ctx)
	if err != nil {
		return nil, apierror.NewAPIError(echo.ErrInternalServerError.Code, err)
	}

	return settings, nil
}
