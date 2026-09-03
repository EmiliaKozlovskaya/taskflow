package tasks_transport_http

import (
	"net/http"
	"time"

	core_kafka "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/kafka"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
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
	//для отправки в Kafka
	requestedAt := time.Now().UTC()

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

	event := core_kafka.NewTaskEvent(requestedAt, core_kafka.TaskDeleted)

	if err := h.publisher.Publish(ctx, event); err != nil {
		log.Warn("failed to publish task event", zap.Error(err))
	}

	responseHandler.NoContentResponse(rw)
}
