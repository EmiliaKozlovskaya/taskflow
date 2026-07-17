package users_service

import (
	"context"
	"fmt"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
)

// в уровне сервиса вся бизнес-логика(валидация )
func (s *UsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	//1. get user by id from repo
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}
	//2. aplly patch to user
	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}
	//3. save patched user in repo
	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user: %w", err)
	}
	return patchedUser, nil
}
