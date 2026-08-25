package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

// DeleteTask   godoc
// @Summary     Удаление задачи
// @Description Удаление существующей в системе задачи по её ID
// @Tags        tasks
// @Param       id path int true "ID удаляемой задачи"
// @Success     204 "Успешное удаление задачи"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server Error"
// @Router      /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	//id удаляемой задачи передается в параметрах пути, поэтому нам надо его достать
	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get taskID path value",
		)
		return
	}
	if err := h.tasksService.DeleteTask(ctx, taskID); err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse(rw)
}
