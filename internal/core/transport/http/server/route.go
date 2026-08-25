package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/middleware"
)

// роут - набор параметров, благодаря которым мультиплексор сможет понять, какой хэндлер нужно вызвать для обработки конкретного http запроса
type Route struct {
	Method     string                            //метод http запроса (GET, POST, PUT, DELETE и т.д.)
	Path       string                            //путь, по которому будет доступен хэндлер (например, /users, /users/{id} и т.д.)
	Handler    http.HandlerFunc                  //хэндлер, который будет обрабатывать http запросы по указанному пути и методу
	Middleware []core_http_middleware.Middleware //теперь у каждого роута есть набор миддлварий
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
