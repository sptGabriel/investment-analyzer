package domain

import "errors"

var (
	ErrFailedDependency    = errors.New("failed dependency")
	ErrMalformedParameters = errors.New("malformed parameters")
	ErrConflict            = errors.New("conflict")
	ErrNotFound            = errors.New("not found")
)
