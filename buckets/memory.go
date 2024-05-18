package buckets

import (
	"fmt"

	"github.com/jjmrocha/oblivion/infra"
)

type inMemoryBucket map[string]any

type InMemoryRepo struct {
	storage map[string]inMemoryBucket
}

func NewInMemoryRepo() *InMemoryRepo {
	repo := InMemoryRepo{
		storage: make(map[string]inMemoryBucket),
	}

	return &repo
}

func (r *InMemoryRepo) CreateBucket(name string) (*Bucket, error) {
	if _, found := r.storage[name]; found {
		reason := fmt.Sprintf("Bucket %v already exists", name)
		return nil, infra.NewErroWithReason(infra.BucketNotFound, reason)
	}

	r.storage[name] = make(inMemoryBucket)

	bucket := Bucket{
		Name: name,
	}

	return &bucket, nil
}

func (r *InMemoryRepo) GetBucket(name string) (*Bucket, error) {
	if _, found := r.storage[name]; !found {
		reason := fmt.Sprintf("bucket %v doesn't exists", name)
		return nil, infra.NewErroWithReason(infra.BucketNotFound, reason)
	}

	bucket := Bucket{
		Name: name,
	}

	return &bucket, nil
}
