package users_postgres_repository

import core_postgres_pool "github.com/Emilia20112005/golang-todoapp/internal/core/repository/postgres/pool"

type UsersRepository struct {
	//мы не хотим тут напрямую зависеть от подключения pgx.Conn, потому что при тестировании это значит
	//что придется устанавливать реальное подключение, а нам это не надо и это усложняет юнит-тест
	//поэтому используем интерфейс чтобы потом в него можно было просто запихнуть заглушку и протестировать
	//core_repository
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
