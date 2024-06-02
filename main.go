package main

import (
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/api"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/storage"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// init
	repository := storage.NewSQLDBRepo("sqlite3", "./test.db")
	defer repository.Close()

	buckectService := bucket.NewBucketService(repository)
	api := api.NewApi(buckectService)
	// setup routing
	mux := http.NewServeMux()
	api.SetRoutes(mux)
	// start
	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
