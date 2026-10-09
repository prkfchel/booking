package bookings_postgres_repository

import (
	"booking/internal/core/domain"
	"time"
)

type BookModel struct {
	ID         int
	EventID    int
	UserID     int
	CreatedAt  time.Time
	Canceled   bool
	CanceledAt *time.Time
}

func domainFromModel(model BookModel) domain.Booking {
	return domain.NewBooking(
		model.ID,
		model.EventID,
		model.UserID,
		model.CreatedAt,
		model.Canceled,
		model.CanceledAt,
	)
}

func domainsFromModels(models []BookModel) []domain.Booking {
	domains := make([]domain.Booking, len(models))

	for i, model := range models {
		domains[i] = domainFromModel(model)
	}

	return domains
}
