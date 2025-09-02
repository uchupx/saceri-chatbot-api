package redis

import (
	"context"
	"time"

	"github.com/uchupx/saceri-chatbot-api/internal/models"
	"github.com/uchupx/saceri-chatbot-api/internal/repository"
	"github.com/uchupx/saceri-chatbot-api/pkg/helper"
)

type SettingCacheRepo struct {
	*CacheRepo
	fallback repository.SettingRepoInterface
}

func (repo *SettingCacheRepo) GetByKey(ctx context.Context, key models.SettingKey) (*models.SettingModel, error) {
	settingKey := models.MakeSettingKey(string(key))
	setting, err := get[models.SettingModel](ctx, *repo.CacheRepo, settingKey)
	if err != nil {
		repo.log.Error(ctx, "Failed to unmarshal cached setting", err, nil)
	}

	if setting == nil {
		return repo.fallback.GetByKey(ctx, key)
	}

	return setting, nil
}

func (repo *SettingCacheRepo) Create(ctx context.Context, data models.SettingModel) (*models.SettingModel, error) {
	settingKey := models.MakeSettingKey(string(data.Key))

	defer func() {
		err := repo.cache.Put(ctx, settingKey, data, 0)
		if err != nil {
			repo.log.Error(ctx, "Failed to put cached setting", err, nil)
		}
	}()

	return repo.fallback.Create(ctx, data)
}
func (repo *SettingCacheRepo) Update(ctx context.Context, data models.SettingModel) (*models.SettingModel, error) {
	settingKey := models.MakeSettingKey(string(data.Key))
	if err := repo.cache.Del(ctx, settingKey); err != nil {
		repo.log.Error(ctx, "Failed to delete cached setting", err, nil)
	}

	defer func() {
		strData, err := helper.StructToString(data)
		if err != nil {
			repo.log.Error(ctx, "Failed to put cached setting", err, nil)
			return
		}

		err = repo.cache.Put(ctx, settingKey, strData, 0)
		if err != nil {
			repo.log.Error(ctx, "Failed to put cached setting", err, nil)
		}
	}()

	return repo.fallback.Update(ctx, data)

}

func (repo *SettingCacheRepo) GetAllSettings(ctx context.Context) ([]models.SettingModel, error) {
	settingKey := string("all_settings")
	setting, err := get[[]models.SettingModel](ctx, *repo.CacheRepo, settingKey)
	if err != nil {
		repo.log.Error(ctx, "Failed to unmarshal cached setting", err, nil)
	}

	if setting == nil {
		s, err := repo.fallback.GetAllSettings(ctx)
		if err != nil {
			return nil, err
		}

		strData, err := helper.StructToString(s)
		if err != nil {
			repo.log.Error(ctx, "failed to marshal struct", err, nil)
			return nil, err
		}

		err = repo.cache.Put(ctx, settingKey, strData, time.Hour*1)
		if err != nil {
			repo.log.Error(ctx, "Failed to cache setting", err, map[string]interface{}{"key": settingKey})
			return s, nil
		}

		setting = &s
	}

	return *setting, nil
}

func NewSettingCacheRepo(repo *CacheRepo, fallback repository.SettingRepoInterface) *SettingCacheRepo {
	return &SettingCacheRepo{
		CacheRepo: repo,
		fallback:  fallback,
	}
}
