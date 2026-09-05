package users_postgres_repository

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, version, name, phone_number, mail, created_at FROM booking.users
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"select users: %w",
			err,
		)
	}

	defer rows.Close()

	var models []UserModel

	for rows.Next() {
		var model UserModel
		err = rows.Scan(
			&model.ID,
			&model.Version,
			&model.Name,
			&model.PhoneNumber,
			&model.Mail,
			&model.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan error: %w",
				err,
			)
		}

		models = append(models, model)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	domains := DomainsFromModels(models)

	return domains, nil
}
