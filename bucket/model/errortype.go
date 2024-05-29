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
	InvalidKey
	// Value related
	MissingField
	InvalidField
	UnknownField
	// Request related
	BadRequestPaylod
	// Schema related
	InvalidBucketName
	SchemaMissing
	InvalidSchema
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
	InvalidKey: {
		statusCode: http.StatusBadRequest,
		template:   "Invalid key %v: %v",
	},
	MissingField: {
		statusCode: http.StatusUnprocessableEntity,
		template:   "Missing field: %v",
	},
	InvalidField: {
		statusCode: http.StatusUnprocessableEntity,
		template:   "Invalid value on field %v",
	},
	UnknownField: {
		statusCode: http.StatusUnprocessableEntity,
		template:   "Unknown field: %v",
	},
	InvalidBucketName: {
		statusCode: http.StatusBadRequest,
		template:   "Invalid bucket name",
	},
	SchemaMissing: {
		statusCode: http.StatusBadRequest,
		template:   "Schema must contain at least one field",
	},
	InvalidSchema: {
		statusCode: http.StatusBadRequest,
		template:   "Invalid definition for field %v",
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
