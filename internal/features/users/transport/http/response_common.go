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
	Mail        *string
	CreatedAt   time.Time
}

func UserDTOFromDomain(domain domain.User) UserResponseDTO {
	return UserResponseDTO{
		ID:          domain.ID,
		Version:     domain.Version,
		Name:        domain.Name,
		PhoneNumber: domain.PhoneNumber,
		Mail:        domain.Mail,
		CreatedAt:   domain.CreatedAt,
	}
}

func UsersDTOFromDomains(domains []domain.User) []UserResponseDTO {
	response := make([]UserResponseDTO, len(domains))
	for i, domain := range domains {
		dto := UserDTOFromDomain(domain)
		response[i] = dto
	}

	return response
}
