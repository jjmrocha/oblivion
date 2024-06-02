package main

import (
	"log"
	"net/http"

	"github.com/jjmrocha/oblivion/api"
	"github.com/jjmrocha/oblivion/bucket"
	"github.com/jjmrocha/oblivion/repo"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// init
	repo := repo.New("sqlite3", "./test.db")
	defer repo.Close()

	buckectService := bucket.NewBucketService(repo)
	api := api.NewApi(buckectService)
	// setup routing
	mux := http.NewServeMux()
	api.SetRoutes(mux)
	// start
	log.Println("Server running")
	log.Fatal(http.ListenAndServe(":9090", mux))
}
