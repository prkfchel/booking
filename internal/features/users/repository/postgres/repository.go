package users_postgres_repository

import core_postgres_pool "booking/internal/core/repository/postgres/pool"

type UsersRepository struct {
	Pool core_postgres_pool.Pool
}

func NewUsersRepository(pool core_postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{
		Pool: pool,
	}
}
