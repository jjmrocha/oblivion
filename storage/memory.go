package storage

import (
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
)

type object map[string]any

type bucketDef struct {
	schema []model.Field
	keys   map[string]object
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

func (r *InMemoryRepo) Close() {
	clear(r.storage)
}

func (r *InMemoryRepo) GetAllBuckets() ([]string, error) {
	bucketList := make([]string, 0)

	for name := range r.storage {
		bucketList = append(bucketList, name)
	}

	return bucketList, nil
}

func (r *InMemoryRepo) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	if _, found := r.storage[name]; found {
		return nil, apperror.New(model.BucketAlreadyExits, name)
	}

	bucketDef := bucketDef{
		keys:   make(map[string]object),
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

func (r *InMemoryRepo) Read(bucket *model.Bucket, key string) (map[string]any, error) {
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

func (r *InMemoryRepo) Store(bucket *model.Bucket, key string, value map[string]any) error {
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

func (r *InMemoryRepo) FindKeys(bucket *model.Bucket, query map[string][]any) ([]string, error) {
	bucketDef, found := r.storage[bucket.Name]
	if !found {
		return nil, apperror.New(model.BucketNotFound, bucket.Name)
	}

	keys := make([]string, 0)

	for key, obj := range bucketDef.keys {
		if matches(obj, query) {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func matches(object map[string]any, query map[string][]any) bool {
	for field, criteria := range query {
		value := object[field]
		if match := matchesOne(value, criteria); !match {
			return false
		}
	}

	return true
}

func matchesOne(value any, criteria []any) bool {
	for _, requested := range criteria {
		if value == requested {
			return true
		}
	}

	return false
}
