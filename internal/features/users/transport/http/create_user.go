package users_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type CreateUserRequestDTO struct {
	Name        string  `json:"name" validate:"required,min=3,max=100" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15" example:"+79998887766"`
	Mail        *string `json:"mail"  validate:"omitempty,min=9,max=100" example:"user@mail.com"`
}
type CreateUserResponseDTO UserResponseDTO

// CreateUser 	 godoc
// @Summary      Создать пользователя
// @Description  Создать нового пользователя в системе
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequestDTO true "CreateUser тело запроса"
// @Success      201 {object} CreateUserResponseDTO "Успешно созданный пользователь"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users [post]
func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	var requestDTO CreateUserRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &requestDTO); err != nil {
		responseHandler.ErrorResponse(err, "decode request")
		return
	}

	user := fromDTOToDomain(CreateUserRequestDTO(requestDTO))

	createdUser, err := h.Service.CreateUser(user, ctx)

	if err != nil {
		responseHandler.ErrorResponse(err, "create user")
		return
	}

	response := CreateUserResponseDTO(DTOFromDomain(createdUser))

	responseHandler.JSONResponse(http.StatusCreated, response)

}

func fromDTOToDomain(dto CreateUserRequestDTO) domain.User {
	return domain.NewUninitializedUser(
		dto.Name,
		dto.PhoneNumber,
		dto.Mail,
	)
}
