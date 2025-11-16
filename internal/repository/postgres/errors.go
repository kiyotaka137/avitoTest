package postgres

import "fmt"

type NotFoundError struct{ what string }

func (e NotFoundError) Error() string { return fmt.Sprintf("%s not found", e.what) }
func ErrNotFound(what string) error   { return NotFoundError{what: what} }
