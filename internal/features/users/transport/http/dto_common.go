package users_transport_http

import "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"

// Инфа о созданном ресурсе (вынесли в отдельный файл так как юзаем и в create_users, и get_users и get_user )
type UserDTOResponse struct {
	ID          int     `json:"id"             example:"10"`
	Version     int     `json:"version"        example:"3"`
	FullName    string  `json:"full_name"      example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number"   example:"+79998887766"`
}

// dtoFromDomain — это переводчик «изнутри наружу».
// Он нужен для безопасности и чистоты ответа. Функция берет готового пользователя
// из доменного слоя (после того, как сервис сохранил его в базу и выдал ему реальный ID)
// и упаковывает только разрешенные поля в CreateUserResponse.
// Это гарантирует, что клиент получит в JSON только то, что ему положено увидеть.
func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users)) //создаем слайс UserDTOResponse такой же длины как переданный массив users

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user) //на каждом шаге отдаем дто полученное из домена user
	}
	return usersDTO
}
