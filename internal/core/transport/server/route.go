package core_http_server

import "net/http"

// роут - набор параметров, благодаря которым мультиплексор сможет понять, какой хэндлер нужно вызвать для обработки конкретного http запроса
type Route struct {
	Method  string           //метод http запроса (GET, POST, PUT, DELETE и т.д.)
	Path    string           //путь, по которому будет доступен хэндлер (например, /users, /users/{id} и т.д.)
	Handler http.HandlerFunc //хэндлер, который будет обрабатывать http запросы по указанному пути и методу
}

func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
