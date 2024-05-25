package storage

import "github.com/jjmrocha/oblivion/bucket/model"

type Repository interface {
	GetAllBuckets() ([]*model.Bucket, error)
	CreateBucket(string, []model.Field) (*model.Bucket, error)
	GetBucket(string) (*model.Bucket, error)
	DropBucket(string) error
	Store(*model.Bucket, string, any) error
	Read(*model.Bucket, string) (any, error)
	Delete(*model.Bucket, string) error
	FindKeys(*model.Bucket, map[string][]any) ([]string, error)
}
