package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
	core_errors "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *TasksRepository) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.tasks (title, description, completed, created_at, completed_at, author_user_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, version, title, description, completed, created_at, completed_at, author_user_id;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	//здесь мы должны понимать когда ошибка 500 Internal Server а когда данного автора чей id нам пришел
	//просто нет в бд, то есть это должна быть ошибка 404 Not Found
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" { //код foreign_key_violation
			return domain.Task{}, fmt.Errorf(
				"author user with id='%d' : %w",
				task.AuthorUserID,
				core_errors.ErrNotFound,
			)
		}

		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	//нужно преобразовать теперь taskModel в доменную сущность taskDomain
	taskDomain := TaskDomainFromModel(taskModel)
	return taskDomain, nil
}
