package bucket

import (
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
	"github.com/jjmrocha/oblivion/storage"
)

type BucketService struct {
	repository storage.Repository
}

func NewBucketService(repo storage.Repository) *BucketService {
	service := BucketService{
		repository: repo,
	}
	return &service
}

func (s *BucketService) BucketList() ([]string, error) {
	bucketList, err := s.repository.GetAllBuckets()
	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	bucketNames := make([]string, 0, len(bucketList))

	for _, bucket := range bucketList {
		bucketNames = append(bucketNames, bucket.Name)
	}

	return bucketNames, nil
}

func (s *BucketService) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	bucket, err := s.repository.GetBucket(name)
	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket != nil {
		return nil, apperror.New(model.BucketAlreadyExits, name)
	}

	return s.repository.CreateBucket(name, schema)
}

func (s *BucketService) GetBucket(name string) (*model.Bucket, error) {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, apperror.New(model.BucketNotFound, name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(name string) error {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(model.BucketNotFound, name)
	}

	return s.repository.DropBucket(name)
}

func (s *BucketService) GetValue(name string, key string) (any, error) {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, apperror.New(model.BucketNotFound, name)
	}

	return s.repository.Read(bucket, key)
}

func (s *BucketService) PutValue(name string, key string, value any) error {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(model.BucketNotFound, name)
	}

	return s.repository.Store(bucket, key, value)
}

func (s *BucketService) DeleteValue(name string, key string) error {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(model.BucketNotFound, name)
	}

	return s.repository.Delete(bucket, key)
}