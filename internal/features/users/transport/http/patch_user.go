package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	//для каждого поля нужно рассматривать 3 сценария:
	//1. Поле в JSON не передано
	//2. Поле в JSON передано со значением - меняем на это значение
	//3. Поле в JSON передано как null - меняем в бд на NULL

	/*FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"` -- проблема такого задания в том что мы никак не отличим
	1 и 3 сценарии потому что оба будут преобразовываться в пустую строку, решение -- введём свой тип данных NullableString*/

	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"Максим Максимович"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+71112223344"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be  NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must start with '+'")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser     godoc
// @Summary     Изменение пользователя
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `phone_number` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"phone_number":"+71112223344"` - устанавливает новый номер телефона в БД
// @Description 3. **Явно передано значение null**: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Description Ограничения: `full_name` не может быть выставлен как null
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true 							   "ID изменяемого пользователя"
// @Param       request body PatchUserRequest true             "PatchUser тело запроса"
// @Success     200 {object} PatchUserResponse                 "Пользователь успешно изменен"
// @Failure     400 {object} core_http_response.ErrorResponse  "Bad Request"
// @Failure     404 {object} core_http_response.ErrorResponse  "User Not Found"
// @Failure     409 {object} core_http_response.ErrorResponse  "Conflict"
// @Failure     500 {object} core_http_response.ErrorResponse  "Internal Server error"
// @Router      /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	//нужно получить id
	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to get userID path value",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	//struct UserPatch
	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			rw,
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(rw, response, http.StatusOK)
}

// функция которая json из тела входящего запроса -> в описанную ДТО
func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
