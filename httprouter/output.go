package httprouter

import (
	"encoding/json"
	"log"

	"github.com/jjmrocha/oblivion/apperror"
)

type errorPayload struct {
	Status      int    `json:"status"`
	ErrorCode   int    `json:"error-code"`
	Description string `json:"description"`
}

func errorResponse(err error) *Response {
	errorType := apperror.UnexpectedError
	description := err.Error()

	if appErr, ok := err.(*apperror.Error); ok {
		errorType = appErr.ErrorType
		description = appErr.String()
	}

	statusCode := errorType.StatusCode()

	resp := Response{
		Status: statusCode,
		Payload: errorPayload{
			Status:      statusCode,
			ErrorCode:   errorType.ErrorCode(),
			Description: description,
		},
	}

	return &resp
}

func writeResponse(ctx *Context, resp *Response) {
	ctx.Writer.WriteHeader(resp.Status)

	if resp.Payload != nil {
		ctx.Writer.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(ctx.Writer).Encode(resp.Payload)
		if err != nil {
			log.Printf("Error writing payload %v to response\n", resp.Payload)
		}
	}

	log.Printf("%d: %s: %v\n", resp.Status, ctx.fullRequestURI(), ctx.duration())
}
