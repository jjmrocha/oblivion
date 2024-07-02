package api

import (
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
)

type bucketRepresentation struct {
	Name   string        `json:"name"`
	Schema []model.Field `json:"schema"`
}

func createBucketRepresentation(bucket repo.Bucket) *bucketRepresentation {
	rep := bucketRepresentation{
		Name:   bucket.Name(),
		Schema: bucket.Schema(),
	}

	return &rep
}
