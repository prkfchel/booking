package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) GetBook(
	ctx context.Context,
	id int,
) (domain.Book, error) {
	book, err := s.Repository.GetBook(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"failed to get book with id='%d' from repository: %w",
			id,
			err,
		)
	}

	return book, nil
}
