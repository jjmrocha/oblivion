package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/bucket/model"
)

type errorMsg struct {
	Status    int    `json:"status"`
	ErrorCode int    `json:"error-code"`
	Reason    string `json:"description"`
}

func writeJSONErrorResponse(w http.ResponseWriter, err error) {
	errorType := model.UnexpectedError
	reason := err.Error()

	if appErr, ok := err.(*model.AppError); ok {
		errorType = appErr.ErrorType
		reason = appErr.String()
	}

	statusCode := errorType.StatusCode()
	payload := errorMsg{
		Status:    statusCode,
		ErrorCode: errorType.ErrorCode(),
		Reason:    reason,
	}

	writeJSONResponse(w, statusCode, payload)
}

func writeJSONResponse(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)

	if payload != nil {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(payload)
		if err != nil {
			log.Printf("Error writing payload %v to response\n", payload)
		}
	}
}