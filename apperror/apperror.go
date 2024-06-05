package apperror

import (
	"fmt"
)

type AppError struct {
	ErrorType   ErrorType
	Description string
	Args        []any
	Cause       error
}

func New(errorType ErrorType, args ...any) error {
	return WithCause(errorType, nil, args...)
}

func WithCause(errorType ErrorType, cause error, args ...any) error {
	errorDesc := fmt.Sprintf(errorType.Template(), args...)
	err := AppError{
		ErrorType:   errorType,
		Cause:       cause,
		Description: errorDesc,
		Args:        args,
	}
	return &err
}

func (e *AppError) String() string {
	return e.Description
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrorType, e.Description)
}
