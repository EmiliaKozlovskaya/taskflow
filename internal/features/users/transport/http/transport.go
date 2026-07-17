package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
	core_http_server "github.com/Emilia20112005/golang-todoapp/internal/core/transport/server"
)

// Слои приложения не должны напрямую зависеть друг от друга, а должны взаимодействовать через интерфейсы.
// Поэтому слой транспортов не должен напрямую зависеть от доменной сущности User, а должен использовать
// DTO (Data Transfer Object) - объекты, которые предназначены для передачи данных между слоями приложения.
// DTO могут быть простыми структурами данных, содержащими только те поля, которые необходимы для передачи
// данных между слоями приложения. Это позволяет избежать зависимости между слоями приложения и обеспечивает более гибкую архитектуру приложения.
type UsersHTTPHandler struct {
	usersService UsersService
}

// Интерфейс содержит методы. Любая структура, содержащая эти методы, автоматически реализует этот интерфес
// нужен для реализации уровня сервиса (каждый уровень независим друг от друга), уровень транспорта зависит от интерфейса сервиса
type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		patch domain.UserPatch,
	) (domain.User, error)
}

// конструктор - функция, которая создает и возвращает новый экземпляр структуры.
func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

// это значит что в рамках фичи users мы поддерживаем следующие роуты
func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}", //pathValue по ключу id
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
