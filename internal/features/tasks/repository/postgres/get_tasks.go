package tasks_postgres_repository

import (
	"context"
	"fmt"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

func (r *TasksRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT * FROM todoapp.tasks
	%s
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`
	args := []any{limit, offset} //для r.pool.Query, так как мы не знаем точное количество параметров

	//так как в отличие от limit и offset библиотека pgx сама не может определить передали нам
	//userID или нет и в любом случае будет выполнять эту строчку с WHERE, то мы сами контролируем это
	//вот таким способом
	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id=$3")
		args = append(args, userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := r.pool.Query(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close() //чтобы как можно раньше освободить используемое подключение, которое получили из ConnectionPool на время обработки данного sql запроса

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
	if err = rows.Err(); err != nil { //если возникла ошибка при rows.Next()
		return nil, fmt.Errorf("next rows: %w", err)
	}
	// return taskModelsToDomains(taskModels), nil -- не сработает, тк компилятор будет ожидать три значения, сработало бы если мы только значение функции taskModelsToDomains(taskModels) возваращали
	taskDomains := taskModelsToDomains(taskModels)
	return taskDomains, nil
}
