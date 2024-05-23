package storage

import (
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/exception"
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

func (r *InMemoryRepo) GetAllBuckets() ([]*model.Bucket, error) {
	bucketList := make([]*model.Bucket, 0)

	for bucketName := range r.storage {
		bucket := model.Bucket{
			Name: bucketName,
		}

		bucketList = append(bucketList, &bucket)
	}

	return bucketList, nil
}

func (r *InMemoryRepo) CreateBucket(name string) (*model.Bucket, error) {
	if _, found := r.storage[name]; found {
		return nil, exception.NewError(exception.BucketAlreadyExits, name)
	}

	r.storage[name] = make(inMemoryBucket)

	bucket := model.Bucket{
		Name: name,
	}

	return &bucket, nil
}

func (r *InMemoryRepo) GetBucket(name string) (*model.Bucket, error) {
	if _, found := r.storage[name]; !found {
		return nil, nil
	}

	bucket := model.Bucket{
		Name: name,
	}

	return &bucket, nil
}

func (r *InMemoryRepo) DropBucket(name string) error {
	if _, found := r.storage[name]; !found {
		return nil
	}

	delete(r.storage, name)

	return nil
}

func (r *InMemoryRepo) Read(bucket *model.Bucket, key string) (any, error) {
	store, found := r.storage[bucket.Name]
	if !found {
		return nil, exception.NewError(exception.BucketNotFound, bucket.Name)
	}

	value, found := store[key]
	if !found {
		return nil, exception.NewError(exception.KeyNotFound, key, bucket.Name)
	}

	return value, nil
}
