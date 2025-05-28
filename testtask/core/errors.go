package core

import "errors"

var (
	ErrPersonNotFound = errors.New("person not found")
	ErrAPIFailed      = errors.New("get api failed")
)
