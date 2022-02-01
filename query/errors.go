package query

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	token    tokenType
	next     tokenType
	expected []tokenType
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("token of type %q expected on of %q as next token but got %q", e.token, e.expected, e.next)
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}

var ErrSyntaxError = errors.New("syntax error")
var ErrValidation = errors.New("validation error")

var ErrGroupNotClosed = fmt.Errorf("%w: group was not closed", ErrSyntaxError)
