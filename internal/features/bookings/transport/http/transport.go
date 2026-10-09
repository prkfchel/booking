package bookings_http_transport

import (
	"booking/internal/core/domain"
	core_http_server "booking/internal/core/transport/http/server"
	"context"
	"net/http"
)

type BookHTTPHandler struct {
	Service BookService
}

func NewHTTPBookHandler(s BookService) *BookHTTPHandler {
	return &BookHTTPHandler{
		Service: s,
	}
}

type BookService interface {
	CreateBook(
		ctx context.Context,
		booking domain.Booking,
	) (domain.Booking, error)

	GetBook(
		ctx context.Context,
		bookingID int,
	) (domain.Booking, error)

	GetBooks(
		ctx context.Context,
		filter domain.BookingsFilter,
	) ([]domain.Booking, error)

	CancelBook(
		ctx context.Context,
		bookingID int,
	) (domain.Booking, error)
}

func (h *BookHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/books",
			Handler: h.CreateBook,
		},
		{
			Method:  http.MethodGet,
			Path:    "/books/{id}",
			Handler: h.GetBook,
		},
		{
			Method:  http.MethodGet,
			Path:    "/books",
			Handler: h.GetBooks,
		},
		{
			Method:  http.MethodPost,
			Path:    "/books/{id}/cancel",
			Handler: h.CancelBook,
		},
	}
}
