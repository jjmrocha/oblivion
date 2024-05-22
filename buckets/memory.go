package buckets

import (
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
		return nil, exceptions.NewError(exceptions.BucketNotFound, name)
	}

	r.storage[name] = make(inMemoryBucket)

	bucket := Bucket{
		Name: name,
	}

	return &bucket, nil
}

func (r *InMemoryRepo) GetBucket(name string) (*Bucket, error) {
	if _, found := r.storage[name]; !found {
		return nil, exceptions.NewError(exceptions.BucketNotFound, name)
	}

	bucket := Bucket{
		Name: name,
	}

	return &bucket, nil
}

func (r *InMemoryRepo) DropBucket(name string) error {
	if _, found := r.storage[name]; !found {
		return exceptions.NewError(exceptions.BucketNotFound, name)
	}

	delete(r.storage, name)

	return nil
}
