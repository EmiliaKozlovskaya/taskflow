package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string //alias для типа string, чтобы использовать его как тип для версий API (например, "v1", "v2" и т.д.)

var (
	ApiVersionV1 = ApiVersion("v1") //константа для версии API v1
	ApiVersionV2 = ApiVersion("v2") //константа для версии API v2
	ApiVersionV3 = ApiVersion("v3") //константа для версии API v3
)

//здесь будет API version router
//версионирование API чтобы не ломать старые клиенты при добавлении новых фич, а также чтобы можно было поддерживать несколько версий API одновременно
//  /api/v1/tasks  /api/v2/users/{id}/tasks

type APIVersionRouter struct {
	*http.ServeMux //это встраивание ServeMux -> эта структура автоматически имеет те же методы и !реализует интерфейс Handler! (т.к есть метод ServeHTTP (это единсвтвенное требование интерфейса Handler)
	apiVersion     ApiVersion
}

func NewAPIVersionRouter(
	apiVersion ApiVersion,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path) //формируем паттерн для маршрута, который будет включать метод и путь (например, "GET /users", "POST /users/{id}/tasks" и т.д.)

		r.Handle(pattern, route.Handler) //регистрируем маршрут в ServeMux с помощью метода Handle, который принимает паттерн и хэндлер
	}
}

//то есть цепочка будет такой: APIVersionRouter -> ServeMux -> Route -> Handler, где APIVersionRouter будет обрабатывать маршруты для конкретной версии API, ServeMux будет маршрутизировать запросы на основе паттернов, Route будет содержать информацию о методе, пути и хэндлере, а Handler будет обрабатывать конкретный запрос.
//поэтому далее в server.go пишем метод RegisterAPIRouters()
