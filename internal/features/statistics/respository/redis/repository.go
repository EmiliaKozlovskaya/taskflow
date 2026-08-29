package statistics_redis_repository

import "github.com/redis/go-redis/v9"

// ? могу ли я так же назвать структуру ведь у меня уже есть интерфейс statisticsRepository который должен
// реализовывать метод GetTasks а не Get и Set
type StatisticsRepository struct {
	client *redis.Client
}

func NewStatisticsRepository(
	client *redis.Client,
) *StatisticsRepository {
	return &StatisticsRepository{
		client: client,
	}
}
