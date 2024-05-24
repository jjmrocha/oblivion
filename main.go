package main

import (
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/api"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/storage"
)

func main() {
	// init
	repository := storage.NewInMemoryRepo()
	buckectService := bucket.NewBucketService(repository)
	api := api.NewApi(buckectService)
	// setup routing
	mux := http.NewServeMux()
	api.SetRoutes(mux)
	// start
	log.Fatal(http.ListenAndServe(":9090", mux))
}
