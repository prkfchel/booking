package users_http_transport

import (
	"booking/internal/core/domain"
	"time"
)

type UserResponseDTO struct {
	ID          int       `json:"id" example:"10"`
	Version     int       `json:"version" example:"3"`
	Name        string    `json:"name" example:"Ivan Ivanov"`
	PhoneNumber *string   `json:"phone_number" example:"+79998887766"`
	Email       *string   `json:"email" example:"user@mail.com"`
	CreatedAt   time.Time `json:"created_at" example:"2006-10-09T12:55:45.159001Z"`
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
