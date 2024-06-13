package apperror

import (
	"net/http"
)

type ErrorType int

const (
	_ ErrorType = iota
	// Bucket related
	BucketAlreadyExits
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
	InvalidFieldName
	InvalidFieldType
	// Gneric
	UnexpectedError
)

type config struct {
	statusCode int
	template   string
}

var errorTypes = map[ErrorType]config{
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
		template:   "Invalid key %v",
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
		template:   "Invalid bucket name %v",
	},
	SchemaMissing: {
		statusCode: http.StatusBadRequest,
		template:   "Schema must contain at least one field",
	},
	InvalidFieldName: {
		statusCode: http.StatusBadRequest,
		template:   "Invalid field name %v",
	},
	InvalidFieldType: {
		statusCode: http.StatusBadRequest,
		template:   "Invalid field type %v",
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
	return errorTypes[t].statusCode
}

func (t ErrorType) Template() string {
	return errorTypes[t].template
}
