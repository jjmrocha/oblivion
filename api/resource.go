package api

import (
	"encoding/json"
	"net/http"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
)

type Handler struct {
	service *bucket.BucketService
}

func NewHandler(bucketService *bucket.BucketService) *Handler {
	api := Handler{
		service: bucketService,
	}

	return &api
}

func (h *Handler) SetRoutes(mux *http.ServeMux) {
	setBucketRoutes(mux, h)
	setKeyRoutes(mux, h)
}

func setBucketRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /v1/buckets", func(w http.ResponseWriter, req *http.Request) {
		bucketNames, err := h.service.BucketList()
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusCreated, bucketNames)
	})

	mux.HandleFunc("POST /v1/buckets", func(w http.ResponseWriter, req *http.Request) {
		var request repo.Bucket

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			writeErrorResponse(w, apperror.New(apperror.BadRequestPaylod))
			return
		}

		bucket, err := h.service.CreateBucket(request.Name, request.Schema)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusCreated, &bucket)
	})

	mux.HandleFunc("GET /v1/buckets/{bucket}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")

		bucket, err := h.service.GetBucket(bucketName)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusOK, &bucket)
	})

	mux.HandleFunc("DELETE /v1/buckets/{bucket}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")

		err := h.service.DeleteBucket(bucketName)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusNoContent, nil)
	})
}

func setKeyRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /v1/buckets/{bucket}/keys/{key}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		key := req.PathValue("key")

		value, err := h.service.GetValue(bucketName, key)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusOK, value)
	})

	mux.HandleFunc("PUT /v1/buckets/{bucket}/keys/{key}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		key := req.PathValue("key")

		var value model.Object

		err := json.NewDecoder(req.Body).Decode(&value)
		if err != nil {
			writeErrorResponse(w, apperror.New(apperror.BadRequestPaylod))
			return
		}

		err = h.service.PutValue(bucketName, key, value)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusNoContent, nil)
	})

	mux.HandleFunc("DELETE /v1/buckets/{bucket}/keys/{key}", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		key := req.PathValue("key")

		err := h.service.DeleteValue(bucketName, key)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusNoContent, nil)
	})

	mux.HandleFunc("GET /v1/buckets/{bucket}/keys", func(w http.ResponseWriter, req *http.Request) {
		bucketName := req.PathValue("bucket")
		criteria := req.URL.Query()

		keys, err := h.service.Search(bucketName, criteria)
		if err != nil {
			writeErrorResponse(w, err)
			return
		}

		writeResponse(w, http.StatusOK, keys)
	})
}
