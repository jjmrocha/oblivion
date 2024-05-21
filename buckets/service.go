package buckets

import (
	"github.com/jjmrocha/oblivion/exceptions"
)

type BucketService struct {
	Repository Repository
}

func NewBucketService(repo Repository) *BucketService {
	service := BucketService{
		Repository: repo,
	}
	return &service
}

func (s *BucketService) BucketList() ([]string, error) {
	bucketList, err := s.Repository.GetAllBuckets()
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
	_, err := s.Repository.GetBucket(name)
	if err == nil {
		return nil, exceptions.NewError(exceptions.BucketAlreadyExits, name)
	}

	return s.Repository.CreateBucket(name)
}

func (s *BucketService) GetBucket(name string) (*Bucket, error) {
	bucket, err := s.Repository.GetBucket(name)
	if err != nil {
		return nil, exceptions.NewError(exceptions.BucketNotFound, name)
	}

	return bucket, nil
}
