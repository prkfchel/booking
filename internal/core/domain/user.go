package domain

import (
	"fmt"
	"regexp"
	"time"

	core_errors "booking/internal/core/errors"
)

const (
	userNameMinLen  = 3
	userNameMaxLen  = 100
	userEmailMaxLen = 100
	userEmailMinLen = 9
)

var (
	phoneRe = regexp.MustCompile(`^\+[0-9]{9,14}$`)
	mailRe  = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type User struct {
	ID          int
	Version     int
	Name        string
	PhoneNumber *string
	Mail        *string
	CreatedAt   time.Time
}

func NewUser(
	id int,
	version int,
	name string,
	phoneNumber *string,
	mail *string,
	createdAt time.Time,
) User {
	return User{
		ID:          id,
		Version:     version,
		Name:        name,
		PhoneNumber: phoneNumber,
		Mail:        mail,
		CreatedAt:   createdAt,
	}
}

func NewUninitializedUser(
	name string,
	phoneNumber *string,
	mail *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		name,
		phoneNumber,
		mail,
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
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil && !phoneRe.MatchString(*u.PhoneNumber) {
		return fmt.Errorf(
			"invalid phone number format, expected +XXXXXXXXXX (10-15 chars): %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if u.Mail != nil {
		if err := validateMail(*u.Mail); err != nil {
			return err
		}
	}

	return nil
}

func validateMail(mail string) error {
	if l := len([]rune(mail)); l > userEmailMaxLen || l < userEmailMinLen {
		return fmt.Errorf(
			"email len must be between %d and %d: %w",
			userEmailMinLen,
			userEmailMaxLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if !mailRe.MatchString(mail) {
		return fmt.Errorf(
			"invalid mail format: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}

type PatchUser struct {
	Name        Nullable[string]
	PhoneNumber Nullable[string]
	Mail        Nullable[string]
}

func NewPatchUser(
	name Nullable[string],
	phoneNumber Nullable[string],
	mail Nullable[string],
) PatchUser {
	return PatchUser{
		Name:        name,
		PhoneNumber: phoneNumber,
		Mail:        mail,
	}
}

func (p PatchUser) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf(
			"name can't be set and null: %w",
			core_errors.ErrInvalidArgument,
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

	if patch.Mail.Set {
		tmp.Mail = patch.Mail.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("failed to validate applied patch: %w", err)
	}

	u.Name = tmp.Name
	u.PhoneNumber = tmp.PhoneNumber
	u.Mail = tmp.Mail

	return nil
}
