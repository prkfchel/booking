package users_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getUsersResponse []UserResponseDTO

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	limit, offset, err := GetLimitOffsetQueryParams(r)

	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit and offset query params",
		)
		return
	}

	users, err := h.Service.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := getUsersResponse(UsersDTOFromDomains(users))

	responseHandler.JSONResponse(http.StatusOK, response)
}
