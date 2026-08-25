package users_transport_http

import (
	"net/http"

	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

type GetUserResponse UserDTOResponse

// GetUser      godoc
// @Summary     Получение пользователя
// @Description Получение существующего пользователя из базы по его ID
// @Tags        users
// @Produce     json
// @Param       id path int true                              "ID получаемого пользователя"
// @Success     200 {object} GetUserResponse                  "Пользователь успешно найден"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure     404 {object} core_http_response.ErrorResponse "User Not Found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal Server error"
// @Router      /users/{id} [get]
func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	userID, err := core_http_request.GetIntPathValue(r, "id") //приложение поймет что /123 это именно id при помощи роутов в transport.go
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
