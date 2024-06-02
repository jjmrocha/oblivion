package model

import (
	"fmt"
)

type AppError struct {
	ErrorType   ErrorType
	Description string
	Args        []any
	Reason      error
}

func Error(errorType ErrorType, args ...any) error {
	return ErrorWithReason(errorType, nil, args...)
}

func ErrorWithReason(errorType ErrorType, reason error, args ...any) error {
	errorDesc := fmt.Sprintf(errorType.Template(), args...)
	err := AppError{
		ErrorType:   errorType,
		Reason:      reason,
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
