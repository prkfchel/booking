package users_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type getUsersResponseDTO []UserResponseDTO

// GetUsers      godoc
// @Summary      Список пользователей
// @Description  Просмотр списка пользователей с опциональной пагинацией
// @Tags         users
// @Produce      json
// @Param        limit  query int false                        "Размер страницы с пользователями"
// @Param        offset query int false                        "Смещение страницы с пользователями"
// @Success      200 {object} getUsersResponseDTO              "Успешное получение списка пользователей"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users [get]
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

	response := getUsersResponseDTO(DTOsFromDomains(users))

	responseHandler.JSONResponse(http.StatusOK, response)
}
