package ds

import "errors"

var (
	ErrShortURLAlreadyExists = errors.New("short url already exists")
	ErrShortUrlNotFound      = errors.New("short url not found")
)
