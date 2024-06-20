package apperror

import (
	"fmt"
)

type Error struct {
	ErrorType   ErrorType
	Description string
	Cause       error
}

func (e *Error) String() string {
	return e.Description
}

func (e *Error) Error() string {
	errorMsg := fmt.Sprintf("%v: %v", e.ErrorType, e.Description)

	if e.Cause != nil {
		errorMsg += "\n\tcause by " + e.Cause.Error()
	}

	return errorMsg
}
