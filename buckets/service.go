package buckets

import (
	"github.com/jjmrocha/oblivion/exceptions"
)

type BucketService struct {
	repository Repository
}

func NewBucketService(repo Repository) *BucketService {
	service := BucketService{
		repository: repo,
	}
	return &service
}

func (s *BucketService) BucketList() ([]string, error) {
	bucketList, err := s.repository.GetAllBuckets()
	if err != nil {
		return nil, exceptions.NewErroWithReason(exceptions.UnexpectedError, err)
	}

	bucketNames := make([]string, 0, len(bucketList))

	for _, bucket := range bucketList {
		bucketNames = append(bucketNames, bucket.Name)
	}

	return bucketNames, nil
}

func (s *BucketService) CreateBucket(name string) (*Bucket, error) {
	_, err := s.repository.GetBucket(name)
	if err == nil {
		return nil, exceptions.NewError(exceptions.BucketAlreadyExits, name)
	}

	return s.repository.CreateBucket(name)
}

func (s *BucketService) GetBucket(name string) (*Bucket, error) {
	bucket, err := s.repository.GetBucket(name)
	if err != nil {
		return nil, exceptions.NewError(exceptions.BucketNotFound, name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(name string) error {
	_, err := s.repository.GetBucket(name)
	if err != nil {
		return exceptions.NewError(exceptions.BucketNotFound, name)
	}

	return s.repository.DropBucket(name)
}
