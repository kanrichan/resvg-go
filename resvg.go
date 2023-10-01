package resvg

import "errors"

//go:generate go run internal/gen/gen.go

var (
	ErrWorkerIsBeingUsed = errors.New("worker is being used")
	ErrPointerIsNil      = errors.New("pointer is nil")
)
