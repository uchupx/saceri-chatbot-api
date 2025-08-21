package redis

import (
	"context"

	"github.com/uchupx/saceri-chatbot-api/internal/database"
	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
)

type SettingCacheRepo struct {
	cache    *database.Cache
	fallback repository.SettingRepoInterface
}

func (repo *SettingCacheRepo) GetSettingByKey(ctx context.Context, key models.SettingKey) (*models.SettingModel, error) {
	settingKey := string(key)
	setting, err := repo.cache.Get(ctx, settingKey)
	if err != nil {
		return nil, err
	}

	if setting == "" {
		return repo.fallback.GetByKey(key)
	}

	return &data, nil
}

func NewSettingCacheRepo(cache *database.Cache) *SettingCacheRepo {
	return &SettingCacheRepo{cache: cache}
}
