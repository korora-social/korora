package daoerrors

import "errors"

var (
	NotFound        error = errors.New("Record not found")
	ErrNotSupported error = errors.New("Unsupported database engine")
)
