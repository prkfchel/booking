package events_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type EventResponseDTO struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	TotalTickets int
	CreatedAt    time.Time
	Address      string
	StartsAt     time.Time
	OwnerUserID  int
	Canceled     bool
	CanceledAt   *time.Time
}

func DTOFromDomain(domain domain.Event) EventResponseDTO {
	return EventResponseDTO{
		ID:           domain.ID,
		Version:      domain.Version,
		Title:        domain.Title,
		Description:  domain.Description,
		TotalTickets: domain.TotalTickets,
		CreatedAt:    domain.CreatedAt,
		Address:      domain.Address,
		StartsAt:     domain.StartsAt,
		OwnerUserID:  domain.OwnerUserID,
		Canceled:     domain.Canceled,
		CanceledAt:   domain.CanceledAt,
	}
}

func DTOFromDomains(domains []domain.Event) []EventResponseDTO {
	response := make([]EventResponseDTO, len(domains))

	for i, domain := range domains {
		response[i] = DTOFromDomain(domain)
	}

	return response
}
