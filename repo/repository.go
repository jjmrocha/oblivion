package repo

import (
	"context"

	"github.com/jjmrocha/oblivion/model"
)

type Repository interface {
	Close()
	BucketNames(ctx context.Context) ([]string, error)
	NewBucket(ctx context.Context, name string, schema []model.Field) (Bucket, error)
	GetBucket(ctx context.Context, name string) (Bucket, error)
	DropBucket(ctx context.Context, name string) error
}

type Bucket interface {
	Name() string
	Schema() []model.Field
	Store(ctx context.Context, key string, value model.Object) error
	Read(ctx context.Context, key string) (model.Object, error)
	Delete(ctx context.Context, key string) error
	Keys(ctx context.Context, criteria model.Criteria) ([]string, error)
}
