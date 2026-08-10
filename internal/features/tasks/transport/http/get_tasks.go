package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	//1. user_id
	//2. limit        - получаем в качестве параметров из http запроса
	//3. offset
	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get userID/limit/offset query param",
		)

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get tasks",
		)

		return
	}
	//преобразуем массив tasksDomains в ответные ДТО TaskDTOResponse
	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))

	responseHandler.JSONResponse(
		rw,
		response,
		http.StatusOK,
	)

}

// возвращает опц limit, опц offset и возм ошибку
func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const ( //ключи лучше выносить в константы
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return userID, limit, offset, nil
}
