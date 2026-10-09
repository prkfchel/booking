package users_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type UserResponseDTO struct {
	ID          int
	Version     int
	Name        string
	PhoneNumber *string
	Email       *string
	CreatedAt   time.Time
}

func DTOFromDomain(domain domain.User) UserResponseDTO {
	return UserResponseDTO{
		ID:          domain.ID,
		Version:     domain.Version,
		Name:        domain.Name,
		PhoneNumber: domain.PhoneNumber,
		Email:       domain.Email,
		CreatedAt:   domain.CreatedAt,
	}
}

func DTOsFromDomains(domains []domain.User) []UserResponseDTO {
	response := make([]UserResponseDTO, len(domains))
	for i, domain := range domains {
		dto := DTOFromDomain(domain)
		response[i] = dto
	}

	return response
}
