package domain

import "time"

var (
	UninitializedID         = -1
	UninitializedVersion    = -1
	UninitializedCreatedAt  time.Time
	UninitializedCanceled   = false
	UninitializedCanceledAt *time.Time
)
