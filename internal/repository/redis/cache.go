package redis

import (
	"context"
	"encoding/json"

	"github.com/uchupx/saceri-chatbot-api/internal/database"
	"github.com/uchupx/saceri-chatbot-api/pkg/apilog"
)

type CacheRepo struct {
	cache *database.Cache
	log   *apilog.ApiLog
}

func get[T any](ctx context.Context, repo CacheRepo, key string) (*T, error) {
	result, err := repo.cache.Get(ctx, key)
	if err != nil {
		//repo.log.log.Errorf("Failed to get key %s from cache: %v", key, err)
		return nil, err
	}

	if result == "" {
		return nil, nil
	}

	var data T

	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		//repo.log.Errorf("Failed to unmarshal cached data: %v", err)
		return nil, err
	}

	return &data, nil
}

func NewCacheRepo(cache *database.Cache, log *apilog.ApiLog) *CacheRepo {
	return &CacheRepo{
		cache: cache,
		log:   log,
	}
}
