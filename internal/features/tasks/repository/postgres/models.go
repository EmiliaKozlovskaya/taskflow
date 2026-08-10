package tasks_postgres_repository

import (
	"time"

	"github.com/Emilia20112005/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func TaskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt, //вот здесь происходит неправильное преобразование времени, хоть у нас стоит UTC в бд, когда мы записываем ответ, в ответ идет время системы откуда запущено приложение
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskModelsToDomains(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))
	for i, model := range taskModels {
		domains[i] = TaskDomainFromModel(model)
	}
	return domains
}
