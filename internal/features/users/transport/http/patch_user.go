package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
	core_logger "github.com/Emilia20112005/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/types"
	core_http_utils "github.com/Emilia20112005/golang-todoapp/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	//для каждого поля нужно рассматривать 3 сценария:
	//1. Поле в JSON не передано
	//2. Поле в JSON передано со значением - меняем на это значение
	//3. Поле в JSON передано как null - меняем в бд на NULL

	/*FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"` -- проблема такого задания в том что мы никак не отличим
	1 и 3 сценарии потому что оба будут преобразовываться в пустую строку, решение -- введём свой тип данных NullableString*/

	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
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

// функция которая json из тела входящего запроса -> в описанную ДТО
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponder(log)

	//нужно получить id
	userID, err := core_http_utils.GetIntPathValue(r, "id")
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

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
