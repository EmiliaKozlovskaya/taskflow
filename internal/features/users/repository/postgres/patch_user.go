package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
	core_errors "github.com/Emilia20112005/golang-todoapp/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	//вспоминаем про оптимистичные блокировки,т.е. version тоже увеличиваем
	query := `
	UPDATE todoapp.users
	SET
		full_name=$1,	
		phone_number=$2,
		version=version+1	
	WHERE id=$3 AND version=$4	
	RETURNING
		id,
		version,
		full_name,
		phone_number;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		id,
		user.Version,
	)
	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		//снова может быть несколько вариантов ошибок
		if errors.Is(err, pgx.ErrNoRows) {
			//либо юзера с id не существует, либо существует, но version за это время уже изменилась
			//с данным sql запросом никак эти две ситуации не разделить, но будем опираться на то, что
			//прежде чем патчить юзера мы использовали GetUser, то есть пользователя мы до этого получили
			//значит он существует(существовал), в любом случае за время работы либо поменялась версия
			//либо пользователя в целом удалили, значит -> конфликт
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf(
			"scan error: %w",
			err,
		)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
