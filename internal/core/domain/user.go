package domain

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"

	coreerrors "booking/internal/core/errors"
)

const (
	userNameMinLen  = 3
	userNameMaxLen  = 100
	userEmailMaxLen = 100
)

var phoneRe = regexp.MustCompile(`^\+[0-9]{9,14}$`)

type User struct {
	ID          int
	Version     int
	Name        string
	PhoneNumber *string
	Email       *string
	CreatedAt   time.Time
}

func NewUser(
	id int,
	version int,
	name string,
	phoneNumber *string,
	email *string,
	createdAt time.Time,
) User {
	return User{
		ID:          id,
		Version:     version,
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
		CreatedAt:   createdAt,
	}
}

func NewUninitializedUser(
	name string,
	phoneNumber *string,
	email *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		name,
		phoneNumber,
		email,
		UninitializedCreatedAt,
	)
}

func (u *User) Validate() error {
	nameLen := len([]rune(u.Name))
	if nameLen < userNameMinLen || nameLen > userNameMaxLen {
		return fmt.Errorf(
			"name len must be between %d and %d: %d: %w",
			userNameMinLen,
			userNameMaxLen,
			nameLen,
			coreerrors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil && !phoneRe.MatchString(*u.PhoneNumber) {
		return fmt.Errorf(
			"invalid phone number format, expected +XXXXXXXXXX (10-15 chars): %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if u.Email != nil {
		if err := validateEmail(*u.Email); err != nil {
			return err
		}
	}

	return nil
}

func validateEmail(email string) error {
	if len([]rune(email)) > userEmailMaxLen {
		return fmt.Errorf(
			"email len must be at most %d: %w",
			userEmailMaxLen,
			coreerrors.ErrInvalidArgument,
		)
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return fmt.Errorf(
			"invalid email format: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

type PatchUser struct {
	Name        Nullable[string]
	PhoneNumber Nullable[string]
	Email       Nullable[string]
}

func NewPatchUser(
	name Nullable[string],
	phoneNumber Nullable[string],
	email Nullable[string],
) PatchUser {
	return PatchUser{
		Name:        name,
		PhoneNumber: phoneNumber,
		Email:       email,
	}
}

func (p PatchUser) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf(
			"name can't be set and null: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	return nil
}

func (u *User) ApplyPatch(patch PatchUser) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("failed to validate patch: %w", err)
	}

	tmp := *u

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if patch.Email.Set {
		tmp.Email = patch.Email.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("failed to validate applied patch: %w", err)
	}

	u.Name = tmp.Name
	u.PhoneNumber = tmp.PhoneNumber
	u.Email = tmp.Email

	return nil
}
