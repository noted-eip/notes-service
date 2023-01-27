package models

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists or conflicts with existing resource")
	ErrUnknown       = errors.New("unknown error")
)

type ListOptions struct {
	Limit  int64
	Offset int64
}
