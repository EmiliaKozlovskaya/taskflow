package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse []UserDTOResponse

// GetUsers     godoc
// @Summary     Получение всех пользователей
// @Description Получение всех существующих пользователей в системе c опциональной пагинацией
// @Tags        users
// @Produce     json
// @Param       limit query int false                         "Размер страницы с пользователями"
// @Param       offset query int false                        "Смещение страницы с пользователями"
// @Success     200 {object} GetUsersResponse                 "Пользователи успешно получены"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal Server error"
// @Router      /users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)
	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get 'limit'/'offset' query param",
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get users",
		)
		return
	}
	//если ошибок не было, то список доменных сущностей нужно преобразовать в дтошку и отдать клиенту в http ответе
	//опишем дто которое будет уходить клиенту в ответе
	response := GetUsersResponse(usersDTOFromDomains(userDomains)) //преобразуем к типу чтобы поддерживать чистоту контракта

	responseHandler.JSONResponse(rw, response, http.StatusOK)
}

// возвращает опц limit, опц offset и возм ошибку
func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const ( //ключи лучше выносить в константы
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return limit, offset, nil
}
