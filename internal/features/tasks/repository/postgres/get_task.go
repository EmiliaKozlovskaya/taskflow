package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
	core_errors "github.com/Emilia20112005/golang-todoapp/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT *
	FROM todoapp.tasks
	WHERE id=$1;
	`
	row := r.pool.QueryRow(ctx, query, id)
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
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := TaskDomainFromModel(taskModel)
	return taskDomain, nil
}
