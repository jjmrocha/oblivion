package apperror

import (
	"fmt"
)

type Error struct {
	ErrorType   ErrorType
	Description string
	Cause       error
}

func New(errorType ErrorType, args ...any) error {
	return WithCause(errorType, nil, args...)
}

func WithCause(errorType ErrorType, cause error, args ...any) error {
	errorDesc := fmt.Sprintf(errorType.Template(), args...)
	err := Error{
		ErrorType:   errorType,
		Cause:       cause,
		Description: errorDesc,
	}
	return &err
}

func (e *Error) String() string {
	return e.Description
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%v: %v cause by %v", e.ErrorType, e.Description, e.Cause)
	}

	return fmt.Sprintf("%v: %v", e.ErrorType, e.Description)
}
