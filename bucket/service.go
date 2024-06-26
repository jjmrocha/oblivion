package bucket

import (
	"net/url"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
	"github.com/jjmrocha/oblivion/valid"
)

type BucketService struct {
	repo *repo.Repo
}

func NewService(repo *repo.Repo) *BucketService {
	service := BucketService{
		repo: repo,
	}
	return &service
}

func (s *BucketService) BucketList() ([]string, error) {
	bucketList, err := s.repo.BucketNames()
	if err != nil {
		return nil, apperror.UnexpectedError.NewErrorWithCause(err)
	}

	return bucketList, nil
}

func (s *BucketService) CreateBucket(name string, schema []model.Field) (*repo.Bucket, error) {
	bucket, err := s.repo.GetBucket(name)
	if err != nil {
		return nil, apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket != nil {
		return nil, apperror.BucketAlreadyExits.NewError(name)
	}

	return s.repo.NewBucket(name, schema)
}

func (s *BucketService) GetBucket(name string) (*repo.Bucket, error) {
	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return nil, apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket == nil {
		return nil, apperror.BucketNotFound.NewError(name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(name string) error {
	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket == nil {
		return apperror.BucketNotFound.NewError(name)
	}

	return s.repo.DropBucket(name)
}

func (s *BucketService) Value(name string, key string) (model.Object, error) {
	bucket, err := s.repo.GetBucket(name)
	if err != nil {
		return nil, apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket == nil {
		return nil, apperror.BucketNotFound.NewError(name)
	}

	object, err := bucket.Read(key)
	if err != nil {
		return nil, apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if object == nil {
		return nil, apperror.KeyNotFound.NewError(key, name)
	}

	return object, nil
}

func (s *BucketService) SetValue(name string, key string, value model.Object) error {
	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket == nil {
		return apperror.BucketNotFound.NewError(name)
	}

	err = valid.Object(value, bucket.Schema)
	if err != nil {
		return err
	}

	return bucket.Store(key, value)
}

func (s *BucketService) DeleteValue(name string, key string) error {
	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket == nil {
		return apperror.BucketNotFound.NewError(name)
	}

	return bucket.Delete(key)
}

func (s *BucketService) FindKeys(name string, criteria url.Values) ([]string, error) {
	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return nil, apperror.UnexpectedError.NewErrorWithCause(err)
	}

	if bucket == nil {
		return nil, apperror.BucketNotFound.NewError(name)
	}

	if err := valid.Criteria(criteria, bucket.Schema); err != nil {
		return nil, err
	}

	normalized, err := model.Convert(criteria, bucket.Schema)
	if err != nil {
		return nil, err
	}

	return bucket.Keys(normalized)
}
