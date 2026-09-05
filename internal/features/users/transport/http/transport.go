package users_http_transport

import (
	"booking/internal/core/domain"
	core_http_server "booking/internal/core/transport/http/server"
	"context"
	"net/http"
)

type UsersHTTPHandler struct {
	Service UsersService
}

func NewHTTPUsersHandler(s UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		Service: s,
	}
}

type UsersService interface {
	CreateUser(
		user domain.User,
		ctx context.Context,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		userID int,
	) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		userID int,
	) error

	PatchUser(
		ctx context.Context,
		userID int,
		userPatch domain.PatchUser,
	) (domain.User, error)
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: h.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: h.PatchUser,
		},
	}
}
