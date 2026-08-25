package tasks_postgres_repository

import core_postgres_pool "github.com/EmiliaKozlovskaya/golang-todoapp/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(
	//внедряем зависимости структуры через аргументы конструктора
	pool core_postgres_pool.Pool,
) *TasksRepository {
	return &TasksRepository{
		pool: pool,
	}
}
