package users_transport_http

import (
	"net/http"

	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/utils"
)

// DELETE /users/{id} -- удаление также по id

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get userID path value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to delete user",
		)
		return
	}

	//204 No Content
	responseHandler.NoContentResponse(rw)
}
