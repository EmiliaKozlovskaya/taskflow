package users_transport_http

import (
	"net/http"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_logger "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/logger"
	core_http_request "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/transport/http/response"
)

// DTO нужны для http запросов и ответов. Request - что приходит от клиента, Response - что отправляем клиенту.
// объявляем только поля которые будем использовать
type CreateUserRequest struct {
	//id не нужен(сами нагенерим), version не нужен(по умолчанию 1)
	FullName    string  `json:"full_name"    validate:"required,min=3,max=100"                example:"Ivan Ivanov"` //сразу валидируем
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"  example:"+7999887766"` //examples оставляем для красивой документации в swagger
}

// делаем type alias на общий вид http response
type CreateUserResponse UserDTOResponse

// CreateUser   godoc
// @Summary     Создать пользователя
// @Description Метод для обработки http запроса на создание нового пользователя в системе.
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       request body CreateUserRequest true "CreateUser тело запроса (обязательный формат входящей ДТО)"
// @Success     201 {object} CreateUserResponse "Успешно созданный пользователь возвращается в данном ответном ДТО"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	log.Debug("Invoke CreateUser handler") //вызов CreateUser хэндлера

	responseHandler := core_http_response.NewHTTPResponder(log)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		//1. log
		//2. status_code
		//3. http response json -> error
		//для этого добавили новый метод в http_response_handler.go
		responseHandler.ErrorResponse(rw, err, "failed to decode and validate http request")

		return
	}

	/*if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("Error occured")
	}
	//данный обработчик никак не валидирует приходящего юзера, а так как нам придется это делать часто, то снова выносим в отдельный блок core_http_request

	rw.WriteHeader(http.StatusOK)*/

	userDomain := domainFromDTO(request)
	//вызываем уровень сервиса
	userDomain, err := h.usersService.CreateUser(ctx, userDomain) //здесь управление переходит слою 2.сервиса (потом -> 3.бд) и этот уровень (транспорт) просто ждет ответ
	if err != nil {
		responseHandler.ErrorResponse(rw, err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(rw, response, http.StatusCreated)

}

// domainFromDTO — это переводчик «снаружи внутрь».
// Он берет «сырые» данные, которые клиент прислал по сети в формате DTO,
// и превращает их в чистый объект бизнес-логики (domain.User).
// При этом мы сразу используем специальный конструктор, чтобы задать
// начальные системные поля (например, пустой ID и версию), так как клиент не должен их присылать сам.
func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
