package api

import (
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
)

type externalBucket struct {
	Name   string        `json:"name"`
	Schema []model.Field `json:"schema"`
}

func createExternalBucket(bucket repo.Bucket) *externalBucket {
	rep := externalBucket{
		Name:   bucket.Name(),
		Schema: bucket.Schema(),
	}

	return &rep
}
