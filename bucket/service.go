package bucket

import (
	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/exception"
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
		return nil, exception.NewErroWithReason(exception.UnexpectedError, err)
	}

	bucketNames := make([]string, 0, len(bucketList))

	for _, bucket := range bucketList {
		bucketNames = append(bucketNames, bucket.Name)
	}

	return bucketNames, nil
}

func (s *BucketService) CreateBucket(name string) (*model.Bucket, error) {
	_, err := s.repository.GetBucket(name)
	if err == nil {
		return nil, exception.NewError(exception.BucketAlreadyExits, name)
	}

	return s.repository.CreateBucket(name)
}

func (s *BucketService) GetBucket(name string) (*model.Bucket, error) {
	bucket, err := s.repository.GetBucket(name)
	if err != nil {
		return nil, exception.NewError(exception.BucketNotFound, name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(name string) error {
	_, err := s.repository.GetBucket(name)
	if err != nil {
		return exception.NewError(exception.BucketNotFound, name)
	}

	return s.repository.DropBucket(name)
}
