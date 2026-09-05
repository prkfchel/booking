package users_postgres_repository

import (
	"booking/internal/core/domain"
	"time"
)

type UserModel struct {
	ID          int
	Version     int
	Name        string
	PhoneNumber *string
	Mail        *string
	CreatedAt   time.Time
}

func DomainFromModel(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.Name,
		model.PhoneNumber,
		model.Mail,
		model.CreatedAt,
	)
}

func DomainsFromModels(models []UserModel) []domain.User {
	domains := make([]domain.User, len(models))

	for i, model := range models {
		domain := DomainFromModel(model)
		domains[i] = domain
	}

	return domains
}
