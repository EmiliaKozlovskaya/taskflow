package core_http_response

//Нужен для того, чтобы обрабатывать http ответы и ошибки в одном месте, чтобы не дублировать код в каждом http хендлере.
//Например, если в любом месте приложения произойдет паника, мы хотим отловить ее и вернуть клиенту корректный http ответ с ошибкой,
//а также залогировать эту ошибку в логах. Поэтому создаем отдельный пакет core_http_response, который будет содержать методы для обработки http ответов и ошибок.
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/errors"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

// переделаем структуру относительно исходного кода (так как передавая rw в полях создаем race condition) убираем поле и передаем в качестве аргумента в методы где rw нужен
// также переименуем в просто Responder, потому что это не хэндлер
type HTTPResponder struct {
	log *core_logger.Logger
}

// в конструкторе передаем логгер, чтобы использовать его в методах для логирования ошибок и паник
func NewHTTPResponder(log *core_logger.Logger) *HTTPResponder {
	return &HTTPResponder{log: log}
}

func (h *HTTPResponder) JSONResponse(
	rw http.ResponseWriter,
	responseBody any,
	statusCode int,
) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	if err := json.NewEncoder(rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponder) NoContentResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponder) ErrorResponse(rw http.ResponseWriter, err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound): //например если какой-то пользователь настойчиво пытается получить GET /users/123, а его не сущ, то у нас бы весь лог файл был в ошибках если бы логали на уровне Error, но ведь глобально это не ошибка на уровне приложения, а просто юзер дибил
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(rw, statusCode, err, msg)
}

// мы хотим, чтобы при панике в любом месте приложения, мы могли отловить ее и вернуть клиенту корректный http ответ с ошибкой, а также залогировать эту ошибку в логах. Поэтому создаем метод PanicResponse, который будет вызываться в middleware Panic, когда произойдет паника. В этом методе мы логируем ошибку и возвращаем клиенту http ответ с кодом 500 и сообщением об ошибке.
func (h *HTTPResponder) PanicResponse(rw http.ResponseWriter, p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("Unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err)) //логируем ошибку с сообщением и самим объектом паники, чтобы в логах была информация о том, что произошло и где именно произошла паника

	h.errorResponse(rw, statusCode, err, msg)
}

func (h *HTTPResponder) errorResponse(
	rw http.ResponseWriter,
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	}
	h.JSONResponse(rw, response, statusCode) //т.е это обычный json response но с особым response body с ошибкой

}
