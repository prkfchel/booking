package domain

import (
	"fmt"
	"time"

	coreerrors "booking/internal/core/errors"
)

const (
	eventTitleMinLen       = 2
	eventTitleMaxLen       = 500
	eventDescriptionMinLen = 2
	eventAddressMinLen     = 5
)

type Event struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	TotalTickets int
	CreatedAt    time.Time
	Address      string
	StartsAt     time.Time
	OwnerUserID  int
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
	address string,
	startsAt time.Time,
	ownerUserID int,
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
		Address:      address,
		StartsAt:     startsAt,
		OwnerUserID:  ownerUserID,
		Canceled:     canceled,
		CanceledAt:   canceledAt,
	}
}

func NewUninitializedEvent(
	title string,
	description *string,
	totalTickets int,
	address string,
	startsAt time.Time,
	ownerUserID int,
) Event {
	return NewEvent(
		UninitializedID,
		UninitializedVersion,
		title,
		description,
		totalTickets,
		UninitializedCreatedAt,
		address,
		startsAt,
		ownerUserID,
		UninitializedCanceled,
		UninitializedCanceledAt,
	)
}

func (e *Event) Validate() error {
	titleLen := len([]rune(e.Title))
	if titleLen < eventTitleMinLen || titleLen > eventTitleMaxLen {
		return fmt.Errorf(
			"title len must be between %d and %d: %d: %w",
			eventTitleMinLen,
			eventTitleMaxLen,
			titleLen,
			coreerrors.ErrInvalidArgument,
		)
	}

	if e.Description != nil {
		descriptionLen := len([]rune(*e.Description))
		if descriptionLen < eventDescriptionMinLen {
			return fmt.Errorf(
				"description len must be at least %d: %d: %w",
				eventDescriptionMinLen,
				descriptionLen,
				coreerrors.ErrInvalidArgument,
			)
		}
	}

	if e.TotalTickets < 1 {
		return fmt.Errorf(
			"total tickets must be positive number: %d: %w",
			e.TotalTickets,
			coreerrors.ErrInvalidArgument,
		)
	}

	addressLen := len([]rune(e.Address))
	if addressLen < eventAddressMinLen {
		return fmt.Errorf(
			"address len must be at least %d: %d: %w",
			eventAddressMinLen,
			addressLen,
			coreerrors.ErrInvalidArgument,
		)
	}

	if e.StartsAt.IsZero() {
		return fmt.Errorf(
			"event start time must be set: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if e.OwnerUserID < 1 {
		return fmt.Errorf(
			"owner user id must be positive number: %d: %w",
			e.OwnerUserID,
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

type PatchEvent struct {
	Title        Nullable[string]
	Description  Nullable[string]
	TotalTickets Nullable[int]
	Address      Nullable[string]
	StartsAt     Nullable[time.Time]
}

func NewPatchEvent(
	title Nullable[string],
	description Nullable[string],
	totalTickets Nullable[int],
	address Nullable[string],
	startsAt Nullable[time.Time],
) PatchEvent {
	return PatchEvent{
		Title:        title,
		Description:  description,
		TotalTickets: totalTickets,
		Address:      address,
		StartsAt:     startsAt,
	}
}

func (p PatchEvent) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"title can't be set and null: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if p.TotalTickets.Set && p.TotalTickets.Value == nil {
		return fmt.Errorf(
			"total tickets can't be set and null: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if p.Address.Set && p.Address.Value == nil {
		return fmt.Errorf(
			"address can't be set and null: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if p.StartsAt.Set && p.StartsAt.Value == nil {
		return fmt.Errorf(
			"start time can't be set and null: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func (e *Event) ApplyPatch(patch PatchEvent) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("failed to validate patch: %w", err)
	}

	tmp := *e

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.TotalTickets.Set {
		tmp.TotalTickets = *patch.TotalTickets.Value
	}

	if patch.Address.Set {
		tmp.Address = *patch.Address.Value
	}

	if patch.StartsAt.Set {
		tmp.StartsAt = *patch.StartsAt.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("failed to validate applied patch: %w", err)
	}

	e.Title = tmp.Title
	e.Description = tmp.Description
	e.TotalTickets = tmp.TotalTickets
	e.Address = tmp.Address
	e.StartsAt = tmp.StartsAt

	return nil
}

type EventsFilter struct {
	Pagination

	Title       *string
	Address     *string
	DateFrom    *time.Time
	DateTo      *time.Time
	OwnerUserID *int
}

func NewEventsFilter(
	title *string,
	address *string,
	dateFrom *time.Time,
	dateTo *time.Time,
	ownerUserID *int,
	limit *int,
	offset *int,
) EventsFilter {
	return EventsFilter{
		Pagination:  NewPagination(limit, offset),
		Title:       title,
		Address:     address,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		OwnerUserID: ownerUserID,
	}
}

func (f EventsFilter) Validate() error {
	if err := f.Pagination.Validate(); err != nil {
		return err
	}

	if err := validateOptionalID("owner user id", f.OwnerUserID); err != nil {
		return err
	}

	return validateDateRange(f.DateFrom, f.DateTo)
}
