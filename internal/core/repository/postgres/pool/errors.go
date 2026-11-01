package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows")
	ErrViolatesForeignKey = errors.New("violates foreign key")
	ErrUnknown            = errors.New("unknown")
	ErrViolatesUnique     = errors.New("violates unique")
	ErrViolatesCheck      = errors.New("violates check")
)
