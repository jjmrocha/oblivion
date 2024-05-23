package storage

import "github.com/jjmrocha/oblivion/bucket/model"

type Repository interface {
	GetAllBuckets() ([]*model.Bucket, error)
	CreateBucket(string) (*model.Bucket, error)
	GetBucket(string) (*model.Bucket, error)
	DropBucket(string) error
	//Store(*Bucket, string, any) error
	//Read(*Bucket, string) (any, error)
	//Delete(*Bucket, string) error
}
