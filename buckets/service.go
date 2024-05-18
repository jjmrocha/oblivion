package buckets

import (
	"fmt"

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

func (s *BucketService) CreateBucket(name string) (*Bucket, error) {
	_, err := s.Repository.GetBucket(name)
	if err == nil {
		reason := fmt.Sprintf("Bucket %v already exists", name)
		return nil, exceptions.NewErroWithReason(exceptions.BucketAlreadyExits, reason)
	}

	return s.Repository.CreateBucket(name)
}

func (s *BucketService) BucketList() ([]string, error) {
	bucketList, err := s.Repository.GetAllBuckets()
	if err != nil {
		return nil, exceptions.NewErroWithReason(exceptions.UnexpectedError, err.Error())
	}

	bucketNames := make([]string, 0, len(bucketList))

	for _, bucket := range bucketList {
		bucketNames = append(bucketNames, bucket.Name)
	}

	return bucketNames, nil
}
