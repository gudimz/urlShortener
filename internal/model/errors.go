package model

import "errors"

var (
	ErrNotFound     = errors.New("url not found")
	ErrAlreadyExist = errors.New("identifier already exist")
)
