package domain

import (
	core_errors "booking/internal/core/errors"
	"fmt"
	"regexp"
	"time"
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
	if nameLen < 3 || nameLen > 100 {
		return fmt.Errorf(
			"invalid name length: %d: %w",
			nameLen,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))

		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid phone number length: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid phone number format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if u.Mail != nil {
		mailLen := len([]rune(*u.Mail))

		if mailLen < 9 || mailLen > 100 {
			return fmt.Errorf(
				"invalid mail length: %d: %w",
				mailLen,
				core_errors.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

		if !re.MatchString(*u.Mail) {
			return fmt.Errorf(
				"invalid mail format: %w",
				core_errors.ErrInvalidArgument,
			)
		}
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

func (p *PatchUser) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf(
			"name cant be set and be null: %w",
			core_errors.ErrInvalidArgument,
		)
	}
	return nil
}

func (user *User) AppyPatch(patch PatchUser) error {
	err := patch.Validate()
	if err != nil {
		return fmt.Errorf(
			"failed to validate patch: %w",
			err,
		)
	}

	tmp := *user

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}

	if patch.Name.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if patch.Mail.Set {
		tmp.Mail = patch.Mail.Value
	}

	if err = tmp.Validate(); err != nil {
		return fmt.Errorf(
			"failed to validate applied patch: %w",
			err,
		)
	}

	user.Name = tmp.Name
	user.PhoneNumber = tmp.PhoneNumber
	user.Mail = tmp.Mail

	return nil
}
