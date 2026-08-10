package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

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
