package events_postgres_repository

import core_postgres_pool "booking/internal/core/repository/postgres/pool"

type EventsRepository struct {
	Pool core_postgres_pool.Pool
}

func NewEventsRepository(pool core_postgres_pool.Pool) *EventsRepository {
	return &EventsRepository{
		Pool: pool,
	}
}
