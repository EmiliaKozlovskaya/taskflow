package tasks_service

import (
	"context"
	"fmt"

	"github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	//1. task.Validate() (если вдруг в будущем у нас уровень транспорта будет не http, а напримерт kafka - брокер сообщений, которым занимались не мы и там что-то не то, то оно и дальше пойдет с ошибкой)
	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("validate task domain: %w", err)
	}
	//2. newTask := repo.Save(task)
	task, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}
	//3. return newTask
	return task, nil
}
