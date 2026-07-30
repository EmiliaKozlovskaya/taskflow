package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Emilia20112005/golang-todoapp/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1;
	`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 { //это будет значить, что наш запрос ничего не удалил => в нашей таблице не было юзера с данным id
		return fmt.Errorf("user with id='%d': '%w'", id, core_errors.ErrNotFound)
	}

	return nil
}
