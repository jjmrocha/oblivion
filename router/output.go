package router

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

func writeErrorResponse(ctx *Context, err error) {
	errorType := apperror.UnexpectedError
	description := err.Error()

	if appErr, ok := err.(*apperror.Error); ok {
		errorType = appErr.ErrorType
		description = appErr.String()
	}

	statusCode := errorType.StatusCode()

	ctx.response = &response{
		status: statusCode,
		payload: errorPayload{
			Status:      statusCode,
			ErrorCode:   errorType.ErrorCode(),
			Description: description,
		},
	}

	log.Printf("ERROR => %s %s => %v", ctx.Request.Method, ctx.Request.RequestURI, err.Error())

	writeResponse(ctx)
}

func writeResponse(ctx *Context) {
	ctx.Writer.WriteHeader(ctx.response.status)

	if ctx.response.payload != nil {
		ctx.Writer.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(ctx.Writer).Encode(ctx.response.payload)
		if err != nil {
			log.Printf("Error writing payload %v to response\n", ctx.response.payload)
		}
	}

	log.Printf("%d => %s %s\n", ctx.response.status, ctx.Request.Method, ctx.Request.RequestURI)
}
