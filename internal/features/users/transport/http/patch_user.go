package users_http_transport

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	core_http_type "booking/internal/core/transport/http/type"
	"fmt"
	"net/http"
	"regexp"
)

type PatchUserRequestDTO struct {
	Name        core_http_type.Nullable[string] `json:"name"`
	PhoneNumber core_http_type.Nullable[string] `json:"phone_number"`
	Mail        core_http_type.Nullable[string] `json:"mail"`
}

type PatchUserResponseDTO UserResponseDTO

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getUserIDPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "get id from path value")
		return
	}

	var request PatchUserRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	userPatch := userPatchDomainFromDTO(request)

	user, err := h.Service.PatchUser(ctx, id, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponseDTO(DTOFromDomain(user))

	responseHandler.JSONResponse(http.StatusOK, response)
}

func userPatchDomainFromDTO(request PatchUserRequestDTO) domain.PatchUser {
	return domain.NewPatchUser(
		request.Name.ToDomain(),
		request.PhoneNumber.ToDomain(),
		request.Mail.ToDomain(),
	)
}

func (r *PatchUserRequestDTO) Validate() error {
	if r.Name.Set {
		if r.Name.Value == nil {
			return fmt.Errorf(
				"name cant be set null: %w",
				core_errors.ErrInvalidArgument,
			)
		}
		nameLen := len([]rune(*r.Name.Value))
		if nameLen < 3 || nameLen > 100 {
			return fmt.Errorf(
				"invalid name length: %d: %w",
				nameLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if r.PhoneNumber.Set && r.PhoneNumber.Value != nil {
		phoneNumberLen := len([]rune(*r.PhoneNumber.Value))

		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid phone number length: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*r.PhoneNumber.Value) {
			return fmt.Errorf(
				"invalid phone number format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if r.Mail.Set && r.Mail.Value != nil {
		mailLen := len([]rune(*r.Mail.Value))

		if mailLen < 9 || mailLen > 100 {
			return fmt.Errorf(
				"invalid mail length: %d: %w",
				mailLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

		if !re.MatchString(*r.Mail.Value) {
			return fmt.Errorf(
				"invalid mail format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
