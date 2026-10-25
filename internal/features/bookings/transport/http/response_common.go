package bookings_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type BookResponseDTO struct {
	ID         int        `json:"id" example:"10"`
	EventID    int        `json:"event_id" example:"1"`
	UserID     int        `json:"user_id" example:"1"`
	CreatedAt  time.Time  `json:"created_at" example:"2006-10-09T12:55:45.159001Z"`
	CanceledAt *time.Time `json:"canceled_at" example:"null"`
}

func DTOFromDomain(domain domain.Booking) BookResponseDTO {
	return BookResponseDTO{
		ID:         domain.ID,
		EventID:    domain.EventID,
		UserID:     domain.UserID,
		CreatedAt:  domain.CreatedAt,
		CanceledAt: domain.CanceledAt,
	}
}

func DTOsFromDomains(domains []domain.Booking) []BookResponseDTO {
	response := make([]BookResponseDTO, len(domains))

	for i, domain := range domains {
		response[i] = DTOFromDomain(domain)
	}

	return response
}
