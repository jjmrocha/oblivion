package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/apperror"
)

type errorMsg struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"error-code"`
	Reason    string `json:"description"`
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	errorType := apperror.UnexpectedError
	reason := err.Error()

	if appErr, ok := err.(*apperror.Error); ok {
		errorType = appErr.ErrorType
		reason = appErr.String()
	}

	statusCode := errorType.StatusCode()
	payload := errorMsg{
		Status:    statusCode,
		ErrorCode: errorType.ErrorCode(),
		Reason:    reason,
	}

	writeResponse(w, statusCode, payload)
}

func writeResponse(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)

	if payload != nil {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			log.Printf("Error writing payload %v to response\n", payload)
		}
	}
}
