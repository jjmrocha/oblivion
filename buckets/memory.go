package buckets

import (
	"fmt"

	"github.com/jjmrocha/oblivion/exceptions"
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

func (r *InMemoryRepo) GetAllBuckets() ([]*Bucket, error) {
	bucketList := make([]*Bucket, 0)

	for bucketName := range r.storage {
		bucket := Bucket{
			Name: bucketName,
		}

		bucketList = append(bucketList, &bucket)
	}

	return bucketList, nil
}

func (r *InMemoryRepo) CreateBucket(name string) (*Bucket, error) {
	if _, found := r.storage[name]; found {
		reason := fmt.Sprintf("Bucket %v already exists", name)
		return nil, exceptions.NewErroWithReason(exceptions.BucketNotFound, reason)
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
		return nil, exceptions.NewErroWithReason(exceptions.BucketNotFound, reason)
	}

	bucket := Bucket{
		Name: name,
	}

	return &bucket, nil
}
