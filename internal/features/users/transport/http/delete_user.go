package users_http_transport

import (
	core_logger "booking/internal/core/logger"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

// DeleteUser    godoc
// @Summary      Удаление пользователя
// @Description  Удаление существующего в системе пользователя по его ID
// @Tags         users
// @Param        id  path string true                          "ID удаляемого пользователя" Format(uuid)
// @Success      204                                           "Успешное удаление пользователя"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "User not found"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getUserIDPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "filed to get user id")
		return
	}

	err = h.Service.DeleteUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.NoContentResponse()
}
