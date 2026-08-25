package users_service

import (
	"context"
	"fmt"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	//тут имплементируем логику на уровне сервиса
	//1. user.Validate()
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}
	//2. repo.Save(user) (для начала создадим метод CreateUser у interface UsersRepository в service.go)
	user, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}
	//3. return user
	return user, nil
}
