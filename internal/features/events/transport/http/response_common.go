package events_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type eventResponseDTO struct {
	ID           int
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

func DTOFromDomain(domain domain.Event) eventResponseDTO {
	return eventResponseDTO{
		ID:           domain.ID,
		Title:        domain.Title,
		Description:  domain.Description,
		TotalTickets: domain.TotalTickets,
		CreatedAt:    domain.CreatedAt,
		EventAddress: domain.EventAddress,
		EventTime:    domain.EventTime,
		OwnersUserId: domain.OwnersUserId,
		Canceled:     domain.Canceled,
		CanceledAt:   domain.CanceledAt,
	}
}

func DTOFromDomains(domains []domain.Event) []eventResponseDTO {
	response := make([]eventResponseDTO, len(domains))

	for i, domain := range domains {
		response[i] = DTOFromDomain(domain)
	}

	return response
}
