package errors

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrUnableToResolve = errors.New("unable to resolve")
	ErrLoaderNotFound  = errors.New("loader not found")
)

func WrongType(expected, actual interface{}) error {
	return fmt.Errorf("wrong type: expected %T, got %T", expected, actual)
}
