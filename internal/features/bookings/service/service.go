package bookings_service

import (
	"booking/internal/core/domain"
	"context"
)

type BookService struct {
	Repository BookRepository
}

func NewBookService(r BookRepository) *BookService {
	return &BookService{
		Repository: r,
	}
}

type BookRepository interface {
	CreateBook(
		ctx context.Context,
		book domain.Booking,
	) (domain.Booking, error)

	GetBook(
		ctx context.Context,
		bookID int,
	) (domain.Booking, error)

	GetBooks(
		ctx context.Context,
		filter domain.BookingsFilter,
	) ([]domain.Booking, error)

	CancelBook(
		ctx context.Context,
		bookID int,
	) (domain.Booking, error)
}
