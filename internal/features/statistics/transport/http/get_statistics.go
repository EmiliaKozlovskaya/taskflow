package statisctics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`          //если вдруг задач вообще нет то и процент нам неоткуда считать
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"` //не *time.Duration чтобы можно было не в наносекундах а в удобном представлении вернуть ответ
}

func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	//нужно прочитать query parameters опциональные GET /statisctics?user_id={user_id}&from={from}&to={to}
	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(rw, err, "failed to get userID/from/to query params")
		return
	}
	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get statistics",
		)
		return
	}
	//преобразуем доменную сущность статистики в ответную ДТО
	response := DTOFromDomain(statistics)

	responseHandler.JSONResponse(rw, response, http.StatusOK)
}

func DTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string //чтобы не было потом паники разыменования (dereference) nil pointer
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String() //нужна переменная посредник пмшт нельзя взять указатель на результат функции
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const ( //ключи лучше выносить в константы
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}
	return userID, from, to, nil
}
