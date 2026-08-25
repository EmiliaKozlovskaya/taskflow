package statisctics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_http_server "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatiscticsService
}

type StatiscticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statiscticsService StatiscticsService, //внедрение зависимостей, благодаря которому под интерфейс statisticsService можем засунуть всё что хотим
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statiscticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/statistics",
			Handler: h.GetStatistics,
		},
	}
}
