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
	EventAddress string
	EventTime    time.Time
	OwnersUserId int
	Canceled     bool
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
		model.EventAddress,
		model.EventTime,
		model.OwnersUserId,
		model.Canceled,
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
