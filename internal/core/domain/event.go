package domain

import (
	core_errors "booking/internal/core/errors"
	"fmt"
	"time"
)

type Event struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	TotalTickets int
	CreatedAt    time.Time
	EventAddress string
	EventTime    time.Time
	OwnersUserId int
	Canceled     bool
	CanceledAt   *time.Time
}

func NewEvent(
	id int,
	version int,
	title string,
	description *string,
	totalTickets int,
	createdAt time.Time,
	eventAddress string,
	eventTime time.Time,
	ownersUserId int,
	canceled bool,
	canceledAt *time.Time,
) Event {
	return Event{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		TotalTickets: totalTickets,
		CreatedAt:    createdAt,
		EventAddress: eventAddress,
		EventTime:    eventTime,
		OwnersUserId: ownersUserId,
		Canceled:     canceled,
		CanceledAt:   canceledAt,
	}
}

func CreateUninitializedEvent(
	title string,
	description *string,
	totalTickets int,
	eventAddress string,
	eventTime time.Time,
	ownersUserId int,
) Event {
	return NewEvent(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		totalTickets,
		UninitializedCreatedAt,
		eventAddress,
		eventTime,
		ownersUserId,
		UninitializedCanceled,
		UninitializedCanceledAt,
	)
}

func (e *Event) Validate() error {
	titleLen := len([]rune(e.Title))

	if titleLen < 2 || titleLen > 500 {
		return fmt.Errorf(
			"title len must be between 2 and 500: %d: %w",
			titleLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if e.Description != nil {
		descriptionLen := len([]rune(*e.Description))
		if descriptionLen < 2 {
			return fmt.Errorf(
				"description len must be more then 2: %d: %w",
				descriptionLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if e.TotalTickets < 1 {
		return fmt.Errorf(
			"total tickets must be positive number: %d: %w",
			e.TotalTickets,
			core_errors.ErrInvalidArgument,
		)
	}

	eventAddresLen := len([]rune(e.EventAddress))

	if eventAddresLen < 5 {
		return fmt.Errorf(
			"event address len must be more then 5: %d: %w",
			eventAddresLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if e.EventTime.IsZero() {
		return fmt.Errorf(
			"evet time must be: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if e.OwnersUserId < 1 {
		return fmt.Errorf(
			"owners user id must be positive number: %d, %w",
			e.OwnersUserId,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type PatchEvent struct {
	Title        Nullable[string]
	Description  Nullable[string]
	TotalTickets Nullable[int]
	EventAddress Nullable[string]
	EventTime    Nullable[time.Time]
	Canceled     Nullable[bool]
}

func NewPatchEvent(
	title Nullable[string],
	description Nullable[string],
	totalTickets Nullable[int],
	eventAddress Nullable[string],
	eventTime Nullable[time.Time],
	canceled Nullable[bool],
) PatchEvent {
	return PatchEvent{
		Title:        title,
		Description:  description,
		TotalTickets: totalTickets,
		EventAddress: eventAddress,
		EventTime:    eventTime,
		Canceled:     canceled,
	}
}

func (patch PatchEvent) Validate() error {
	if patch.Title.Set && patch.Title.Value == nil {
		return fmt.Errorf(
			"title cant be set and null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if patch.TotalTickets.Set && patch.TotalTickets.Value == nil {
		return fmt.Errorf(
			"total tickets cant be set and null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if patch.EventAddress.Set && patch.EventAddress.Value == nil {
		return fmt.Errorf(
			"event address tickets cant be set and null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if patch.EventTime.Set && patch.EventTime.Value == nil {
		return fmt.Errorf(
			"event time tickets cant be set and null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if patch.Canceled.Set && patch.Canceled.Value == nil {
		return fmt.Errorf(
			"canceled tickets cant be set and null: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

func (event *Event) ApplyPatch(patch PatchEvent) error {
	err := patch.Validate()
	if err != nil {
		return fmt.Errorf(
			"failed to validate patch: %w",
			err,
		)
	}

	tmp := *event

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.TotalTickets.Set {
		tmp.TotalTickets = *patch.TotalTickets.Value
	}

	if patch.EventAddress.Set {
		tmp.EventAddress = *patch.EventAddress.Value
	}

	if patch.EventTime.Set {
		tmp.EventTime = *patch.EventTime.Value
	}

	if patch.Canceled.Set {
		tmp.Canceled = *patch.Canceled.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf(
			"failed to validate applied patch: %w",
			err,
		)
	}

	event.Title = tmp.Title
	event.Description = tmp.Description
	event.TotalTickets = tmp.TotalTickets
	event.EventAddress = tmp.EventAddress
	event.EventTime = tmp.EventTime
	event.Canceled = tmp.Canceled

	return nil
}

type EventsQueryFilter struct {
	Title        *string
	EventAddress *string
	DateFrom     *time.Time
	DateTo       *time.Time
	OwnersUserId *int
	Limit        *int
	Offset       *int
}

func NewEventsQueryFilter(
	title *string,
	eventAddress *string,
	dateFrom *time.Time,
	dateTo *time.Time,
	ownersUserId *int,
	limit *int,
	offset *int,
) EventsQueryFilter {
	return EventsQueryFilter{
		Title:        title,
		EventAddress: eventAddress,
		DateFrom:     dateFrom,
		DateTo:       dateTo,
		OwnersUserId: ownersUserId,
		Limit:        limit,
		Offset:       offset,
	}
}

func (filter EventsQueryFilter) Validate() error {
	if filter.Limit != nil && *filter.Limit < 1 {
		return fmt.Errorf(
			"limit must be positive: %d: %w",
			*filter.Limit,
			core_errors.ErrInvalidArgument,
		)
	}

	if filter.Offset != nil && *filter.Offset < 1 {
		return fmt.Errorf(
			"offset must be positive: %d: %w",
			*filter.Offset,
			core_errors.ErrInvalidArgument,
		)
	}

	if filter.OwnersUserId != nil && *filter.OwnersUserId < 0 {
		return fmt.Errorf(
			"owners user id must be positive: %d: %w",
			*filter.OwnersUserId,
			core_errors.ErrInvalidArgument,
		)
	}

	if filter.DateFrom != nil && filter.DateTo != nil {
		dateFrom := *filter.DateFrom
		dateTo := *filter.DateTo
		if !dateFrom.Before(dateTo) {
			return fmt.Errorf(
				"dateFrom='%d' must be bofare dateTo='%d': %w",
				dateFrom,
				dateTo,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}
