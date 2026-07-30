package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT * 
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1 
	OFFSET $2;
	`
	//если limit/offset = nil то они просто не будут учитываться при выполнении sql запроса
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close() //по правилам библиотеки pgx нужно высвободить подключение

	var userModels []UserModel

	for rows.Next() {
		var userModel UserModel
		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return []domain.User{}, fmt.Errorf("scan error: %w", err)
		}
		//кладем просканированную модель в слайс моделей
		userModels = append(userModels, userModel)
	}
	//после всего цикла также необходимо проверять на возможную ошибку
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w:", err)
	}
	//далее нужно этот слайс моделей преобразовать в слайс доменов (в файле models.go)
	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
