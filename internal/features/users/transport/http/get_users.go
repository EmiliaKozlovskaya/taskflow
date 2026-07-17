package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

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
	limit, err := core_http_utils.GetQueryParams(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	offset, err := core_http_utils.GetQueryParams(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}
	return limit, offset, nil
}
