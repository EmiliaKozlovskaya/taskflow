package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout()) //дочерний контекст тому который передали с таймаутом
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
	`
	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	//Благодаря строчке RETURNING id, version, full_name, phone_number;
	//Postgres выполнил запись и сразу вернул новые системные данные.
	//Мы обязаны прочитать их из row и упаковать обратно в domain.User,
	// чтобы вернуть наверх (в сервис и транспорт) полноценный, официально созданный в системе объект.

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}
	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)
	return userDomain, nil

}
