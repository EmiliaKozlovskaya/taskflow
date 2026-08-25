package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask      godoc
// @Summary     Получение задачи
// @Description Получение существующей задачи из базы по её ID
// @Tags        tasks
// @Produce     json
// @Param       id path int true                              "ID получаемой задачи"
// @Success     200 {object} GetTaskResponse                  "Задача успешно найдена"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task Not Found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal Server error"
// @Router      /tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get taskID path value",
		)
		return
	}
	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(
		rw,
		//taskDomain - не можем вернуть так как Доменная сущность - это бизнес-логика, а транспорт - это API-контракт, они должны быть независимы
		response,
		http.StatusOK,
	)
}
