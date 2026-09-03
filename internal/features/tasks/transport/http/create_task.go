package tasks_transport_http

import (
	"net/http"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_kafka "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/kafka"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

// DTO вход запроса, здесь отражаем что мы хотим получать во входящем запросе при создании задачи
type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100"          example:"Домашнее задание"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000"  example:"Сделать до четверга"`
	AuthorUserID int     `json:"author_user_id" validate:"required"               example:"5"`
}

// DTO ответа
type CreateTaskResponse TaskDTOResponse

// CreateTask   godoc
// @Summary     Создать задачу
// @Description Метод для обработки http запроса на создание новой задачи в системе.
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       request body CreateTaskRequest true "CreateTask тело запроса (обязательный формат входящей ДТО)"
// @Success     201 {object} CreateTaskResponse "Успешно созданная задача возвращается в данном ответном ДТО"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Author not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)
	//для отправки в Kafka
	requestedAt := time.Now().UTC()

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}
	//далее нужно какую-то доменную сущность передать в уровень сервиса, но пока её нет
	//поэтому создаем её в internal/core/domain
	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to create task",
		)
		return
	}

	event := core_kafka.NewTaskEvent(requestedAt, core_kafka.TaskCreated)

	if err := h.publisher.Publish(ctx, event); err != nil {
		log.Warn("failed to publish task event", zap.Error(err))
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain)) //приводим к нашему типу конкретного хэндлера

	responseHandler.JSONResponse(
		rw,
		response,
		http.StatusCreated,
	)
}
