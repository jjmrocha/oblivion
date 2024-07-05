package repo

import (
	"github.com/jjmrocha/oblivion/model"
)

type Repository interface {
	Close()
	BucketNames() ([]string, error)
	NewBucket(name string, schema []model.Field) (Bucket, error)
	GetBucket(name string) (Bucket, error)
	DropBucket(name string) error
}

type Bucket interface {
	Name() string
	Schema() []model.Field
	Store(key string, value model.Object) error
	Read(key string) (model.Object, error)
	Delete(key string) error
	Keys(criteria model.Criteria) ([]string, error)
}
