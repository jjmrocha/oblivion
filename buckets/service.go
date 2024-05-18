package buckets

import (
	"fmt"

	"github.com/jjmrocha/oblivion/infra"
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
		return nil, infra.NewErroWithReason(infra.BucketAlreadyExits, reason)
	}

	return s.Repository.CreateBucket(name)
}
