package tasks_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_kafka "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/kafka"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/types"
	"go.uber.org/zap"
)

type PatchTaskRequest struct {
	//по бизнес-правилам должны уметь менять заголовок/описание/completed
	Title       core_http_types.Nullable[string] `json:"title"        swaggertype:"string" example:"Сделать домашку"`
	Description core_http_types.Nullable[string] `json:"description"  swaggertype:"string" example:"Сделать до четверга"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"    swaggertype:"boolean"`
}

// чтобы функция DecodeAndValidate применила наши кастомные правила валидации
func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			//по бизнес требованиям не может быть Title==null
			return fmt.Errorf("`Title` can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}
	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
			}
		}
	}
	if r.Completed.Set {
		//здесь проверяем только не пришло ли значение null, потому что completed является обязательным полем
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}
	return nil
}

type PatchTaskResponse TaskDTOResponse

// PatchTask    godoc
// @Summary     Изменение задачи
// @Description Изменение информации об уже существующей в системе задачи
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"description":"Сделать до четверга"` - устанавливает новое описание в БД
// @Description 3. **Явно передано значение null**: `"description": null` - очищает поле в БД (set to NULL)
// @Description Ограничения: `title` и `completed` не может быть выставлен как null
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id      path      int             true 		       "ID изменяемой задачи"
// @Param       request body     PatchTaskRequest true             "PatchTask тело запроса"
// @Success     200     {object} PatchTaskResponse                 "Задача успешно изменена"
// @Failure     400     {object} core_http_response.ErrorResponse  "Bad Request"
// @Failure     404     {object} core_http_response.ErrorResponse  "Task Not Found"
// @Failure     409     {object} core_http_response.ErrorResponse  "Conflict"
// @Failure     500     {object} core_http_response.ErrorResponse  "Internal Server error"
// @Router      /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)
	//для отправки в Kafka
	requestedAt := time.Now().UTC()

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(rw, err, "failed to get taskID path value")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil { //передаем по ссылке так как внутри DecodeAndValidateRequest при вызове .Decode автоматически вызывает наш переопределенный метод UnmarshallJSON который пихает тело запроса конкретно в эту структурку нашу (не создавая копии)
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	//далее нужно преобразовать в доменную сущность чтобы отправлять дальше по слоям
	taskPatch := taskPatchDomainFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to patch task",
		)
		return
	}

	event := core_kafka.NewTaskEvent(requestedAt, core_kafka.TaskUpdated)

	if err := h.publisher.Publish(ctx, event); err != nil {
		log.Warn("failed to publish task event", zap.Error(err))
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(rw, response, http.StatusOK)

}

func taskPatchDomainFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
