package models

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists or conflicts with existing resource")
	ErrUnknown       = errors.New("unknown error")
	ErrForbidden     = errors.New("forbidden operation")
)

type ListOptions struct {
	Limit  int32
	Offset int32
}
