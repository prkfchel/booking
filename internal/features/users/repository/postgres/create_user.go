package users_postgres_repository

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
	"time"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO booking.users (name, phone_number, mail, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, name, phone_number, mail, created_at;
	`

	row := r.Pool.QueryRow(ctx, query, user.Name, user.PhoneNumber, user.Mail, time.Now())

	var model UserModel

	err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Name,
		&model.PhoneNumber,
		&model.Mail,
		&model.CreatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := DomainFromModel(model)

	return userDomain, nil
}
