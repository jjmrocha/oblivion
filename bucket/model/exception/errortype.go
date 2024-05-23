package exception

import (
	"net/http"
)

type ErrorType int

const (
	// Bucket related
	BucketAlreadyExits ErrorType = iota + 1
	BucketNotFound
	// Keys related
	KeyNotFound
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
	KeyNotFound: {
		statusCode:  http.StatusNotFound,
		description: "Key %v not found on bucket %v",
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

func (t ErrorType) ErrorCode() int {
	return int(t)
}

func (t ErrorType) StatusCode() int {
	return errorTypeDefMap[t].statusCode
}
