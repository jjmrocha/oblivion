package exception

import (
	"fmt"
)

type AppError struct {
	ErrorType   ErrorType
	Description string
	Reason      error
}

func NewError(errorType ErrorType, args ...any) error {
	return NewErrorWithReason(errorType, nil, args...)
}

func NewErrorWithReason(errorType ErrorType, reason error, args ...any) error {
	errorDesc := fmt.Sprintf(errorTypeDefMap[errorType].description, args...)
	err := AppError{
		ErrorType:   errorType,
		Reason:      reason,
		Description: errorDesc,
	}
	return &err
}

func (e *AppError) String() string {
	return e.Description
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%v: %v", e.ErrorType, e.Description)
}
