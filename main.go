package main

import (
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/api"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/storage"
)

func routes(api *api.Api) *http.ServeMux {
	mux := http.NewServeMux()
	// Bucket operations
	mux.HandleFunc("GET /v1/buckets", api.GetAllBuckets)
	mux.HandleFunc("POST /v1/buckets", api.CreateBucket)
	mux.HandleFunc("GET /v1/buckets/{bucket}", api.GetBucket)
	mux.HandleFunc("DELETE /v1/buckets/{bucket}", api.DeleteBucket)
	// key operations
	mux.HandleFunc("GET /v1/buckets/{bucket}/keys/{key}", api.ReadKey)
	mux.HandleFunc("PUT /v1/buckets/{bucket}/keys/{key}", api.UpdateKey)
	mux.HandleFunc("DELETE /v1/buckets/{bucket}/keys/{key}", api.DeleteKey)
	mux.HandleFunc("GET /v1/buckets/{bucket}/keys", api.Search)
	return mux
}

func main() {
	repository := storage.NewInMemoryRepo()
	buckectService := bucket.NewBucketService(repository)
	api := api.NewApi(buckectService)
	mux := routes(api)
	log.Fatal(http.ListenAndServe(":9090", mux))
}
