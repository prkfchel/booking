package users_http_transport

import (
	"booking/internal/core/domain"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	"net/http"
)

type CreateUserRequestBody struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15"`
	Mail        *string `json:"mail"  validate:"omitempty,min=9,max=100"`
}
type CreateUserResponse UserResponseDTO

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	var requestDTO CreateUserRequestBody

	if err := core_http_request.DecodeAndValidateRequest(r, &requestDTO); err != nil {
		responseHandler.ErrorResponse(err, "decode request")
		return
	}

	domain := fromDTOToDomain(CreateUserRequestBody(requestDTO))

	userDomain, err := h.Service.CreateUser(domain, ctx)

	if err != nil {
		responseHandler.ErrorResponse(err, "create user")
		return
	}

	response := CreateUserResponse(UserDTOFromDomain(userDomain))

	responseHandler.JSONResponse(http.StatusCreated, response)

}

func fromDTOToDomain(dto CreateUserRequestBody) domain.User {
	return domain.NewUninitializedUser(
		dto.Name,
		dto.PhoneNumber,
		dto.Mail,
	)
}
