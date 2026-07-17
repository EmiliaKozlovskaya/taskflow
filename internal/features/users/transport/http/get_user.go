package users_transport_http

import (
	"net/http"

	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	userID, err := core_http_utils.GetIntPathValue(r, "id") //приложение поймет что /123 это именно id при помощи роутов в transport.go
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get userID path value",
		)
		return
	}

	// GET /users/{id} нужно получать id, данную функцию вынесем в core_http_utils
	userDomain, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get user",
		)
		return
	}
	response := GetUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(rw, response, http.StatusOK)
}
