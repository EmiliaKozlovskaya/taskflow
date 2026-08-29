package statistics_service

import (
	"context"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

type StatisticsCache interface {
	Get(
		ctx context.Context,
		key string,
	) (domain.Statistics, bool, error)
	Set(
		ctx context.Context,
		key string,
		statistics domain.Statistics,
		ttl time.Duration,
	) error
}

type StatisticsService struct {
	statisticsRepository StatisticsRepository
	statisticsCache      StatisticsCache
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
	statisticsCache StatisticsCache,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
		statisticsCache:      statisticsCache,
	}
}
