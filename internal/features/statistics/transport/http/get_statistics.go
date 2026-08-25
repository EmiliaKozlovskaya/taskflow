package statisctics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"                  example:"50"`
	TasksCompleted             int      `json:"tasks_completed"                example:"10"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"           example:"20"`    //если вдруг задач вообще нет то и процент нам неоткуда считать
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"  example:"1m30s"` //не *time.Duration чтобы можно было не в наносекундах а в удобном представлении вернуть ответ
}

// GetStatistics     godoc
// @Summary          Получение статистики
// @Description      Получение статистики по задачам с опциональной фильтрацией по user_id и/или временному промежутку
// @Tags             statistics
// @Produce          json
// @Param            user_id query    int    false             "Фильтрация статистики по конкретному пользователю"
// @Param            from    query    string false             "Начало промежутка рассмотрения статистики(включительно), формат: YYYY-MM-DD"
// @Param            to      query    string false             "Конец промежутка рассмотрения статистики(не включительно), формат: YYYY-MM-DD"
// @Success          200     {object} GetStatisticsResponse    "Пользователи успешно получены"
// @Failure          400     {object} core_http_response.ErrorResponse    "Bad Request"
// @Failure          500     {object} core_http_response.ErrorResponse    "Internal Server error"
// @Router          /statistics [get]
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
