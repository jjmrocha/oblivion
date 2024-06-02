package api

import (
	"encoding/json"
	"net/http"

	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/repo"
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
	setBucketRoutes(mux, api)
	setKeyRoutes(mux, api)
}

func setBucketRoutes(mux *http.ServeMux, api *Api) {
	mux.HandleFunc("GET /v1/buckets", func(w http.ResponseWriter, req *http.Request) {
		bucketNames, err := api.bucketService.BucketList()
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusCreated, bucketNames)
	})

	mux.HandleFunc("POST /v1/buckets", func(w http.ResponseWriter, req *http.Request) {
		var request repo.Bucket

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			writeJSONErrorResponse(w, model.Error(model.BadRequestPaylod))
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
}

func setKeyRoutes(mux *http.ServeMux, api *Api) {
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

		if len(key) == 0 || len(key) > 50 {
			writeJSONErrorResponse(w, model.Error(model.InvalidKey, key, "0 < key <= 50"))
			return
		}

		var value map[string]any

		err := json.NewDecoder(req.Body).Decode(&value)
		if err != nil {
			writeJSONErrorResponse(w, model.Error(model.BadRequestPaylod))
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
		bucketName := req.PathValue("bucket")
		query := req.URL.Query()

		keys, err := api.bucketService.Search(bucketName, query)
		if err != nil {
			writeJSONErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusOK, keys)
	})
}
