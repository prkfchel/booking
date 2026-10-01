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
		book domain.Book,
	) (domain.Book, error)

	GetBook(
		ctx context.Context,
		id int,
	) (domain.Book, error)

	GetBooks(
		ctx context.Context,
		filter domain.BooksQueryFilter,
	) ([]domain.Book, error)

	CancelBook(
		ctx context.Context,
		id int,
	) (domain.Book, error)
}
