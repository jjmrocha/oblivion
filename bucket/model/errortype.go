package model

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
	statusCode int
	template   string
}

var errorTypeDefMap = map[ErrorType]errorTypeDef{
	BucketAlreadyExits: {
		statusCode: http.StatusConflict,
		template:   "Bucket %v already exists",
	},
	BucketNotFound: {
		statusCode: http.StatusNotFound,
		template:   "Bucket %v not found",
	},
	KeyNotFound: {
		statusCode: http.StatusNotFound,
		template:   "Key %v not found on bucket %v",
	},
	BadRequestPaylod: {
		statusCode: http.StatusBadRequest,
		template:   "Bad request: Invalid body",
	},
	UnexpectedError: {
		statusCode: http.StatusInternalServerError,
		template:   "Unexpected error",
	},
}

func (t ErrorType) ErrorCode() int {
	return int(t)
}

func (t ErrorType) StatusCode() int {
	return errorTypeDefMap[t].statusCode
}

func (t ErrorType) Template() string {
	return errorTypeDefMap[t].template
}
