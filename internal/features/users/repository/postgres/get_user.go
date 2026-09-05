package users_postgres_repository

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, name, phone_number, mail, created_at FROM booking.users
		WHERE id = $1;
	`

	row := r.Pool.QueryRow(ctx, query, userID)

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
		return domain.User{}, fmt.Errorf(
			"select error: %w",
			err,
		)
	}

	userDomain := DomainFromModel(model)

	return userDomain, nil
}
