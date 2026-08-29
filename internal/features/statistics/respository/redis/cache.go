package statistics_redis_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	"github.com/redis/go-redis/v9"
)

func (r *StatisticsRepository) Get(
	ctx context.Context,
	key string,
) (domain.Statistics, bool, error) { //поле bool добавляем для понимания нашлась ли статистика в redis или нет (found/!found)
	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return domain.Statistics{}, false, nil //в данном случае мы не нашли статистику в кеше, поэтому возвращаем пустую структуру и nil ошибку
	}
	if err != nil {
		return domain.Statistics{}, false, fmt.Errorf("get statistics from redis: %w", err)
	}
	var statistics domain.Statistics
	if err := json.Unmarshal([]byte(value), &statistics); err != nil {
		return domain.Statistics{}, false, fmt.Errorf("unmarshall cached statistics: %w", err)
	}
	return statistics, true, nil
}

func (r *StatisticsRepository) Set(
	ctx context.Context,
	key string,
	statistics domain.Statistics,
	ttl time.Duration,
) error {
	//не сработает со statistics, так как это структура непонятная для redis, поэтому нужно сериализовать в json ([]byte)
	/*if err := r.client.Set(ctx, key, statistics, ttl).Err(); err != nil {
		//...
	}*/
	value, err := json.Marshal(statistics)
	if err != nil {
		return fmt.Errorf("marshall statistics for cache: %w", err)
	}
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("set statistics to redis: %w", err)
	}
	return nil
}
