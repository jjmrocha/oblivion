package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/buckets"
	"github.com/jjmrocha/oblivion/exceptions"
)

type Api struct {
	bucketService *buckets.BucketService
}

func NewApi(bucketService *buckets.BucketService) *Api {
	api := Api{
		bucketService: bucketService,
	}

	return &api
}

func (api *Api) CreateBucket(w http.ResponseWriter, req *http.Request) {
	var request CreateBucketRequest

	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		writeJSONErrorResponse(w, exceptions.NewError(exceptions.BadRequestPaylod))
		return
	}

	bucket, err := api.bucketService.CreateBucket(request.Name)
	if err != nil {
		writeJSONErrorResponse(w, err)
		return
	}

	response := BucketResponse{
		Name: bucket.Name,
	}

	writeJSONResponse(w, http.StatusCreated, &response)
}

func (api *Api) GetAllBuckets(w http.ResponseWriter, req *http.Request) {
	bucketNames, err := api.bucketService.BucketList()
	if err != nil {
		writeJSONErrorResponse(w, err)
		return
	}

	writeJSONResponse(w, http.StatusCreated, bucketNames)
}

func (api *Api) GetBucket(w http.ResponseWriter, req *http.Request) {
	bucketName := req.PathValue("bucket")

	bucket, err := api.bucketService.Repository.GetBucket(bucketName)
	if err != nil {
		writeJSONErrorResponse(w, err)
		return
	}

	paylod := BucketResponse{
		Name: bucket.Name,
	}

	writeJSONResponse(w, http.StatusOK, &paylod)
}

func (api *Api) DeleteBucket(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func (api *Api) UpdateKey(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func (api *Api) ReadKey(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func (api *Api) DeleteKey(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func writeJSONErrorResponse(w http.ResponseWriter, err error) {
	errorType := exceptions.UnexpectedError
	reason := err.Error()

	if appErr, ok := err.(*exceptions.AppError); ok {
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error writing payload %v to response\n", payload)
	}
}
