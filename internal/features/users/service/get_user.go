package users_service

import (
	"context"
	"fmt"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	//валидации никакой не надо, потому что если присылают несуществующий id то просто возваращется 404
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user from repository: %w", err)
	}
	return user, nil
}
