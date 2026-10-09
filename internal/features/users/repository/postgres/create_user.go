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

	queryCreateUser := `
		INSERT INTO booking.users (name, phone_number, mail, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version, name, phone_number, mail, created_at;
	`

	row := r.Pool.QueryRow(ctx, queryCreateUser, user.Name, user.PhoneNumber, user.Email, time.Now())

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

	createdUser := DomainFromModel(model)

	return createdUser, nil
}
