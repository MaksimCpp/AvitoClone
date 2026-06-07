package item

import "errors"

var (
	ErrItemNotFound = errors.New("Item not found.")
	ErrForbidden = errors.New("Forbidden.")
)