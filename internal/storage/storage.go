package storage

import "errors"

var (
	ErrURLNotFound = errors.New("URL Not Found")
	ErrURLExists = errors.New("URL Exists")
)