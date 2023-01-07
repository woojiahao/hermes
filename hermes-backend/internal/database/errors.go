package database

import "errors"

var (
	NotFoundError = errors.New("not found")
	InternalError = errors.New("internal")
	DataError     = errors.New("data error")
)
