package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
)

type Api struct {
	bucketService *bucket.BucketService
}

func NewApi(bucketService *bucket.BucketService) *Api {
	api := Api{
		bucketService: bucketService,
	}

	return &api
}

func (api *Api) SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/buckets", func(w http.ResponseWriter, req *http.Request) {
		bucketNames, err := api.bucketService.BucketList()
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusCreated, bucketNames)
	})

	mux.HandleFunc("POST /v1/buckets", func(w http.ResponseWriter, req *http.Request) {
		var request model.Bucket

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			writeJSONErrorResponse(w, apperror.New(model.BadRequestPaylod))
			return
		}

		err = checkBucketCreation(request.Name, request.Schema)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		bucket, err := api.bucketService.CreateBucket(request.Name, request.Schema)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusCreated, &bucket)
	})

	mux.HandleFunc("GET /v1/buckets/{bucket}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")

		bucket, err := api.bucketService.GetBucket(bucketName)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusOK, &bucket)
	})

	mux.HandleFunc("DELETE /v1/buckets/{bucket}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")

		err := api.bucketService.DeleteBucket(bucketName)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusNoContent, nil)
	})

	// key operations
	mux.HandleFunc("GET /v1/buckets/{bucket}/keys/{key}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		key := req.PathValue("key")

		value, err := api.bucketService.GetValue(bucketName, key)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusOK, value)
	})

	mux.HandleFunc("PUT /v1/buckets/{bucket}/keys/{key}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		key := req.PathValue("key")

		var value map[string]any

		err := json.NewDecoder(req.Body).Decode(&value)
		if err != nil {
			writeJSONErrorResponse(w, apperror.New(model.BadRequestPaylod))
			return
		}

		err = api.bucketService.PutValue(bucketName, key, value)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusNoContent, nil)
	})

	mux.HandleFunc("DELETE /v1/buckets/{bucket}/keys/{key}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		key := req.PathValue("key")

		err := api.bucketService.DeleteValue(bucketName, key)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusNoContent, nil)
	})

	mux.HandleFunc("GET /v1/buckets/{bucket}/keys", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})
}

func writeJSONErrorResponse(w http.ResponseWriter, err error) {
	errorType := model.UnexpectedError
	reason := err.Error()

	if appErr, ok := err.(*apperror.AppError); ok {
		errorType = appErr.ErrorType
		reason = appErr.String()
	}

	statusCode := errorType.StatusCode()
	payload := Error{
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
