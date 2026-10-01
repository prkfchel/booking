package bookings_postgres_repository

import core_postgres_pool "booking/internal/core/repository/postgres/pool"

type BookRepository struct {
	Pool core_postgres_pool.Pool
}

func NewBookRepository(p core_postgres_pool.Pool) *BookRepository {
	return &BookRepository{
		Pool: p,
	}
}
