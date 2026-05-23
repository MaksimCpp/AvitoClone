package item

import "errors"

var (
	ErrNotFound = errors.New("Item not found.")
	ErrForbidden = errors.New("Forbidden.")
)