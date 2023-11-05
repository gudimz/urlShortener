package ds

import "errors"

var (
	ErrShortURLAlreadyExists = errors.New("short url already exists")
	ErrShortURLNotFound      = errors.New("short url not found")
)
