package events_http_transport

import (
	"booking/internal/core/domain"
	core_errors "booking/internal/core/errors"
	core_logger "booking/internal/core/logger"
	core_http_request "booking/internal/core/transport/http/request"
	core_http_response "booking/internal/core/transport/http/response"
	core_http_type "booking/internal/core/transport/http/type"
	"fmt"
	"net/http"
	"time"
)

type patchEventRequestDTO struct {
	Title        core_http_type.Nullable[string]    `json:"title"`
	Description  core_http_type.Nullable[string]    `json:"description"`
	TotalTickets core_http_type.Nullable[int]       `json:"total_tickets"`
	Address      core_http_type.Nullable[string]    `json:"address"`
	StartsAt     core_http_type.Nullable[time.Time] `json:"starts_at"`
}

type patchEventResponseDTO EventResponseDTO

// PatchEvent    godoc
// @Summary      Изменение события
// @Description  Изменение информации об уже существующем в системе событии
// @Description  ### Логика обновления полей (Three-state logic):
// @Description  1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description  2. **Явно передано значение**: `"description": "событие бла-бла-бла"` - устанавливает новое описание в БД
// @Description  3. **Передан null**: `"description": null` - очищает поле в БД (set to NULL)
// @Description  Ограничения: `title` не может быть выставлен как null
// @Description  Ограничения: `total_tickets` не может быть выставлен как null
// @Description  Ограничения: `address` не может быть выставлен как null
// @Description  Ограничения: `starts_at` не может быть выставлен как null
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id      path string           		true       "ID изменяемого события" Format(uuid)
// @Param        request body patchEventRequestDTO 	true       "PatchEvent тело запроса"
// @Success      200 {object} patchEventResponseDTO             "Успешно изменённый пользователь"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "event not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /event/{id} [patch]
func (h *EventsHTTPHandler) PatchEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseWriter(log, w)

	id, err := getEventIDPathValue(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get event id path value")
		return
	}

	var request patchEventRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate")
		return
	}

	eventPatch := PatchEventDoaminFromDTO(request)

	event, err := h.Sevice.PatchEvent(ctx, id, eventPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch event")
		return
	}

	response := patchEventResponseDTO(DTOFromDomain(event))

	responseHandler.JSONResponse(http.StatusOK, response)

}

func PatchEventDoaminFromDTO(request patchEventRequestDTO) domain.PatchEvent {
	return domain.NewPatchEvent(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.TotalTickets.ToDomain(),
		request.Address.ToDomain(),
		request.StartsAt.ToDomain(),
	)
}

func (request *patchEventRequestDTO) Validate() error {
	if request.Title.Set {
		titleLen := len([]rune(*request.Title.Value))

		if titleLen < 2 || titleLen > 500 {
			return fmt.Errorf(
				"title len must be between 2 and 500: %d: %w",
				titleLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if request.Description.Set {
		if request.Description.Value != nil {
			descriptionLen := len([]rune(*request.Description.Value))
			if descriptionLen < 2 {
				return fmt.Errorf(
					"description len must be more then 2: %d: %w",
					descriptionLen,
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}

	if request.TotalTickets.Set {
		if *request.TotalTickets.Value < 1 {
			return fmt.Errorf(
				"total tickets must be positive number: %d: %w",
				*request.TotalTickets.Value,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if request.Address.Set {
		eventAddresLen := len([]rune(*request.Address.Value))

		if eventAddresLen < 5 {
			return fmt.Errorf(
				"event address len must be more then 5: %d: %w",
				eventAddresLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if request.StartsAt.Set {
		if request.StartsAt.Value.IsZero() {
			return fmt.Errorf(
				"evet time must be: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
