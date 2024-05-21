package exceptions

import (
	"fmt"
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

type errorTypeDef struct {
	statusCode  int
	description string
}

var errorTypeDefMap = map[ErrorType]errorTypeDef{
	BucketAlreadyExits: {
		statusCode:  http.StatusConflict,
		description: "Bucket %v already exists",
	},
	BucketNotFound: {
		statusCode:  http.StatusNotFound,
		description: "Bucket %v not found",
	},
	BadRequestPaylod: {
		statusCode:  http.StatusBadRequest,
		description: "Bad request: Invalid body",
	},
	UnexpectedError: {
		statusCode:  http.StatusInternalServerError,
		description: "Unexpected error",
	},
}

type AppError struct {
	ErrorType   ErrorType
	Description string
	Reason      error
}

func NewError(errorType ErrorType, args ...any) error {
	return NewErroWithReason(errorType, nil, args...)
}

func NewErroWithReason(errorType ErrorType, reason error, args ...any) error {
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

func (t ErrorType) ErrorCode() int {
	return int(t)
}

func (t ErrorType) StatusCode() int {
	return errorTypeDefMap[t].statusCode
}
