package statistics_service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_errors "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/errors"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	log := core_logger.FromContext(ctx)
	//проверим чтобы начало промежутка действительно было до конца
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	//1. tasks := get tasks (Redis -> PostgreSQL)
	//2. statistics := calcStatistics(tasks)
	//3. return statistics
	key := statisticsCacheKey(userID, from, to)

	statistics, found, err := s.statisticsCache.Get(ctx, key)
	if err == nil && found { //если статистика нашлась в кеше, то возвращаем её
		log.Debug("statistics found in cache", zap.String("key", key))
		return statistics, nil
	}
	if err != nil { //если ошибка с получкнием статисики, то логируем ошибку и идем в репозиторий
		log.Warn("failed to get statistics from cache, will try to get from repository", zap.Error(err))
	}
	//далее идем в репозиторий, так как статистика не нашлась в кеше
	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get tasks from repository: %w", err)
	}
	statistics = calcStatistics(tasks)

	if err := s.statisticsCache.Set(ctx, key, statistics, time.Minute); err != nil {
		log.Warn("failed to cache statistics", zap.Error(err))
	}
	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil) //значения по умолчанию для явности
	}

	tasksCreated := len(tasks)
	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}
	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100 //чтобы не было целочисленного деления приводим явно к дробному типу

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted >= 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletedRate,
		tasksAverageCompletionTime,
	)

}

// создаёт уникальное имя ключа Redis для конкретного запроса статистики
// ф-ция получает те же пар-ры что и GetStatistics и превращает их в строку-ключ
// вида todoapp:statistics:v1:user:%s:from:%s:to:%s
func statisticsCacheKey(
	userID *int,
	from *time.Time,
	to *time.Time,
) string {
	userPart := "all"
	if userID != nil {
		userPart = strconv.Itoa(*userID)
	}

	fromPart := "all"
	if from != nil {
		fromPart = from.UTC().Format(time.RFC3339)
	}

	toPart := "all"
	if to != nil {
		toPart = to.UTC().Format(time.RFC3339)
	}

	return fmt.Sprintf(
		"todoapp:statistics:v1:user:%s:from:%s:to:%s",
		userPart,
		fromPart,
		toPart,
	)
}
