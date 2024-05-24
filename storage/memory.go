package storage

import (
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
)

type bucketDef struct {
	schema []model.Field
	keys   map[string]any
}

type InMemoryRepo struct {
	storage map[string]bucketDef
}

func NewInMemoryRepo() *InMemoryRepo {
	repo := InMemoryRepo{
		storage: make(map[string]bucketDef),
	}

	return &repo
}

func (r *InMemoryRepo) GetAllBuckets() ([]*model.Bucket, error) {
	bucketList := make([]*model.Bucket, 0)

	for name, bucketDef := range r.storage {
		bucket := model.Bucket{
			Name:   name,
			Schema: bucketDef.schema,
		}

		bucketList = append(bucketList, &bucket)
	}

	return bucketList, nil
}

func (r *InMemoryRepo) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	if _, found := r.storage[name]; found {
		return nil, apperror.New(model.BucketAlreadyExits, name)
	}

	bucketDef := bucketDef{
		keys:   make(map[string]any),
		schema: schema,
	}

	r.storage[name] = bucketDef

	bucket := model.Bucket{
		Name:   name,
		Schema: schema,
	}

	return &bucket, nil
}

func (r *InMemoryRepo) GetBucket(name string) (*model.Bucket, error) {
	bucketDef, found := r.storage[name]
	if !found {
		return nil, nil
	}

	bucket := model.Bucket{
		Name:   name,
		Schema: bucketDef.schema,
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
	bucketDef, found := r.storage[bucket.Name]
	if !found {
		return nil, apperror.New(model.BucketNotFound, bucket.Name)
	}

	value, found := bucketDef.keys[key]
	if !found {
		return nil, apperror.New(model.KeyNotFound, key, bucket.Name)
	}

	return value, nil
}

func (r *InMemoryRepo) Store(bucket *model.Bucket, key string, value any) error {
	bucketDef, found := r.storage[bucket.Name]
	if !found {
		return apperror.New(model.BucketNotFound, bucket.Name)
	}

	bucketDef.keys[key] = value

	return nil
}

func (r *InMemoryRepo) Delete(bucket *model.Bucket, key string) error {
	bucketDef, found := r.storage[bucket.Name]
	if !found {
		return apperror.New(model.BucketNotFound, bucket.Name)
	}

	delete(bucketDef.keys, key)

	return nil
}
