package apperror

import (
	"fmt"

	"github.com/jjmrocha/oblivion/bucket/model"
)

type AppError struct {
	ErrorType   model.ErrorType
	Description string
	Args        []any
	Reason      error
}

func New(errorType model.ErrorType, args ...any) error {
	return WithReason(errorType, nil, args...)
}

func WithReason(errorType model.ErrorType, reason error, args ...any) error {
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
