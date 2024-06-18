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
	errorMsg := fmt.Sprintf("%v: %v", e.ErrorType, e.Description)

	if e.Cause != nil {
		errorMsg += " cause by " + e.Cause.Error()
	}

	return errorMsg
}
