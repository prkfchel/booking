package bookings_service

import (
	"booking/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) CreateBook(
	ctx context.Context,
	book domain.Book,
) (domain.Book, error) {
	if err := book.Validate(); err != nil {
		return domain.Book{}, fmt.Errorf(
			"validation error: %w",
			err,
		)
	}

	createdBook, err := s.Repository.CreateBook(ctx, book)

	if err != nil {
		return domain.Book{}, fmt.Errorf(
			"failed to cvreate book: %w",
			err,
		)
	}

	return createdBook, nil
}
