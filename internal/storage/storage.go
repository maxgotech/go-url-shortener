package storage

import "errors"

var (
	ErrURLNotFound = errors.New("URL Not Found")
	ErrURLNotExists = errors.New("URL Exists")
)