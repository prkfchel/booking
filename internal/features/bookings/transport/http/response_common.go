package bookings_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type BookResponseDTO struct {
	ID         int
	EventID    int
	UserID     int
	CreatedAt  time.Time
	Canceled   bool
	CanceledAt *time.Time
}

func DTOFromDomain(domain domain.Booking) BookResponseDTO {
	return BookResponseDTO{
		ID:         domain.ID,
		EventID:    domain.EventID,
		UserID:     domain.UserID,
		CreatedAt:  domain.CreatedAt,
		Canceled:   domain.Canceled,
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
