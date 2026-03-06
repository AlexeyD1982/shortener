package local_errors

import "errors"

var (
	ErrNotUnique = errors.New("not unique short url")
	ErrNotFound  = errors.New("url not found")
)
