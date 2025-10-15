package redis


import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/southern-martin/zride/backend/services/ai-matching-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/errors"
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) domain.CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) SetDriverLocation(ctx context.Context, driverID uuid.UUID, location domain.Location) error {
	key := fmt.Sprintf("driver:location:%s", driverID.String())
	
	locationJSON, err := json.Marshal(location)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to marshal location", err)
	}

	err = r.client.Set(ctx, key, locationJSON, 5*time.Minute).Err()
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to set driver location in cache", err)
	}

	return nil
}

func (r *CacheRepository) GetDriverLocation(ctx context.Context, driverID uuid.UUID) (*domain.Location, error) {
	key := fmt.Sprintf("driver:location:%s", driverID.String())
	
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NewAppError(errors.CodeNotFound, "driver location not found in cache", err)
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get driver location from cache", err)
	}

	var location domain.Location
	err = json.Unmarshal([]byte(result), &location)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to unmarshal location", err)
	}

	return &location, nil
}

func (r *CacheRepository) SetMatchResult(ctx context.Context, key string, result *domain.MatchResult) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to marshal match result", err)
	}

	err = r.client.Set(ctx, key, resultJSON, 10*time.Minute).Err()
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to set match result in cache", err)
	}

	return nil
}

func (r *CacheRepository) GetMatchResult(ctx context.Context, key string) (*domain.MatchResult, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errors.NewAppError(errors.CodeNotFound, "match result not found in cache", err)
		}
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to get match result from cache", err)
	}

	var matchResult domain.MatchResult
	err = json.Unmarshal([]byte(result), &matchResult)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeInternalError, "failed to unmarshal match result", err)
	}

	return &matchResult, nil
}

func (r *CacheRepository) DeleteMatchResult(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return errors.NewAppError(errors.CodeInternalError, "failed to delete match result from cache", err)
	}

	return nil
}