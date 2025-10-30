package util

import "errors"

var (
	ErrOrgRequired        = errors.New("organization name is required")
	ErrResourceIDRequired = errors.New("resource ID is required")
	ErrModelIDRequired    = errors.New("model ID is required")
	ErrDatasetIDRequired  = errors.New("dataset ID is required")
	ErrSpaceIDRequired    = errors.New("space ID is required")
)
