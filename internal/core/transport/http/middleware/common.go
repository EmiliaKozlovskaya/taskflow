// здесь будут общие middleware для всех http хендлеров
package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

// зачастую middleware это функция, которая оборачивает обработку http запроса, то есть добавляет дополнительную функциональность к существующему обработчику.

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//тут создадим множество Origins, которым мы доверяем
			//В отличие от массива (поиск эл-та O(n)) поиск эл-та в мн-ве O(1)
			allowedOrigins := map[string]struct{}{ //таким образом можем реализовать структуру данных множество - неупорядоченный набор УНИКАЛЬНЫХ значений
				"http://localhost:5050": {},
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)                               // если Origin в списке разрешенных, то добавляем его в заголовки ответа
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS") //какие методы разрешены в запросе
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")       //какие заголовки разрешены в запросе
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// middleware для генерации request_id
func RequestID() Middleware { //чтобы вернуть Middleware, нужно вернуть функцию, соответствующую сигнатуре Middleware(из файла middleware.go), которая принимает http.Handler и возвращает http.Handler
	// Middleware для генерации уникального идентификатора запроса (Request ID) и добавления его в заголовки запроса и ответа.
	return func(next http.Handler) http.Handler { // возвращаем функцию, которая принимает http.Handler и возвращает http.Handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //возвращаем http.HandlerFunc, который реализует интерфейс http.Handler (внутри него мы можем писать код, который будет выполняться при обработке запроса)
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r) // вызываем следующий обработчик в цепочке middleware или сам обработчик запроса
		})
	}
}

// обогощаем логгером который преконфигурирован на автоматическое добавление request_id и url к каждому логу
func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			//преконфигурированный логгер
			l := log.With( //по сути добавляем к логгеру поля (id и url), которые будут автоматически добавляться ко всем логам (на деле же мы создаем новый логгер, который будет использоваться только для этого запроса)
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), l) //создаем новый контекст, в котором будет храниться логгер с добавленными полями (request_id и url), чтобы его можно было использовать в других местах приложения
			next.ServeHTTP(w, r.WithContext(ctx))        // передаем запрос с новым контекстом дальше
		})
	}
}

// Middleware для логирования всех входящих HTTP запросов и исходящих HTTP ответов. Логируем метод, URL, статус код ответа и время обработки запроса.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)
			next.ServeHTTP(rw, r) //здесь происходит вызов самого хэндлера (например CreateUser), который уже имеет и request_id, и обогащенный логгер и отловленную панику

			//здесь хотим еще получать статус код ответа, но напрямую достать его из http.ResponseWriter нельзя, поэтому нужно создать свой ResponseWriter, который будет оборачивать оригинальный и сохранять статус код ответа.
			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

// вызывает следующий middleware в цепочке (или сам обработчик запроса), а также отлавливает и обрабатывает панику (будет отлавливать все дальнейшие паники, но не в RequestID и Logger, так как если эти фундаментальные штуки не сработают, то дальше уже не имеет смысла обрабатывать запрос)
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//получаем из входящего контекста логгер и используем наш w чтобы создать новый httpResponseHandler(/response/handler.go), который будет использоваться для обработки паники и отправки ответа клиенту
			ctx := r.Context()
			//вынесем выделение логгера в отдельный метод в logger.go, чтобы не дублировать код в каждом middleware
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponder(log)

			defer func() { //отлавливаем панику, если она произойдет в следующем обработчике через defer и recover (функция, которая позволяет отловить панику и продолжить выполнение программы)
				if p := recover(); p != nil {
					//если произошла паника, то мы ее отлавливаем и логируем
					//создали для этого отдельный метод PanicResponse в core_http_response/handler.go, который будет логировать ошибку и отправлять клиенту корректный http ответ с ошибкой
					responseHandler.PanicResponse(
						w,
						p,
						"during handle HTTP request occured panic",
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
