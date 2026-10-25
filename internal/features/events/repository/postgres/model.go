package events_postgres_repository

import (
	"booking/internal/core/domain"
	"time"
)

type EventModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	TotalTickets int
	CreatedAt    time.Time
	Address      string
	StartsAt     time.Time
	OwnerUserID  int
	CanceledAt   *time.Time
}

func DomainFromModel(model EventModel) domain.Event {
	return domain.NewEvent(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.TotalTickets,
		model.CreatedAt,
		model.Address,
		model.StartsAt,
		model.OwnerUserID,
		model.CanceledAt,
	)
}
func DomainsFromModels(models []EventModel) []domain.Event {
	domains := make([]domain.Event, len(models))

	for i, model := range models {
		domain := DomainFromModel(model)
		domains[i] = domain
	}

	return domains
}
