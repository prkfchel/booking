package events_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type EventResponseDTO struct {
	ID           int        `json:"id" example:"10"`
	Version      int        `json:"version" example:"3"`
	Title        string     `json:"tite" example:"concert"`
	Description  *string    `json:"description" example:"interesting concert and show"`
	TotalTickets int        `json:"total_tickets" example:"100"`
	CreatedAt    time.Time  `json:"created_at" example:"2006-11-10T19:30:00Z"`
	Address      string     `json:"address" example:"st. Pushkina, b. Kalatushkina"`
	StartsAt     time.Time  `json:"starts_at" example:"2026-11-10T19:30:00Z"`
	OwnerUserID  int        `json:"owner_user_id" example:"1"`
	CanceledAt   *time.Time `json:"canceled_at" example:"null"`
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
