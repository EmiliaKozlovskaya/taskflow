package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT * FROM todoapp.tasks
	`)
	args := []any{}
	conditions := []string{}

	//проверили передан ли userID, если да -> добавляем новое условие фильтрации которое подставим в sql запрос
	if userID != nil {
		args = append(args, userID)
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)))
	}

	if from != nil {
		args = append(args, from)
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)))
	}

	if to != nil {
		args = append(args, to)
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args))) //строго < чтобы не дублировались данные
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC; ")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel
	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
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
			return nil, fmt.Errorf("scan tasks: %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	taskDomains := taskModelsToDomains(taskModels)
	return taskDomains, nil
}
