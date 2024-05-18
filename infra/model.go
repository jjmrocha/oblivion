package infra

import (
	"net/http"
)

type ErrorType int

const (
	// Bucket related
	BucketAlreadyExits ErrorType = iota + 1
	BucketNotFound
	// Keys related
	// Request related
	BadRequestPaylod
	// Gneric
	UnexpectedError
)

func (t ErrorType) String() string {
	switch t {
	case BucketAlreadyExits:
		return "Bucket already exists"
	case BucketNotFound:
		return "Bucket not found"
	case BadRequestPaylod:
		return "Bad request: Invalid body"
	}

	return "Unexpected error"
}

func (t ErrorType) StatusCode() int {
	switch t {
	case BucketAlreadyExits:
		return http.StatusConflict
	case BucketNotFound:
		return http.StatusNotFound
	case BadRequestPaylod:
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func (t ErrorType) ErrorCode() int {
	return int(t)
}

type AppError struct {
	ErrorType ErrorType
	Reason    *string
}

func NewError(errorType ErrorType) error {
	err := AppError{
		ErrorType: errorType,
	}
	return &err
}

func NewErroWithReason(errorType ErrorType, reason string) error {
	err := AppError{
		ErrorType: errorType,
		Reason:    &reason,
	}
	return &err
}

func (e *AppError) String() string {
	if e.Reason != nil {
		return *e.Reason
	}

	return e.ErrorType.String()
}

func (e *AppError) Error() string {
	return e.String()
}
