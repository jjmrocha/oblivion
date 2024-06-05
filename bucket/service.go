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
	bucketList, err := s.repo.GetAllBuckets()
	if err != nil {
		return nil, apperror.WithCause(apperror.UnexpectedError, err)
	}

	return bucketList, nil
}

func (s *BucketService) CreateBucket(name string, schema []model.Field) (*repo.Bucket, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	if err := valid.Schema(schema); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)
	if err != nil {
		return nil, apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket != nil {
		return nil, apperror.New(apperror.BucketAlreadyExits, name)
	}

	return s.repo.CreateBucket(name, schema)
}

func (s *BucketService) GetBucket(name string) (*repo.Bucket, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return nil, apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, apperror.New(apperror.BucketNotFound, name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(name string) error {
	if err := valid.BucketName(name); err != nil {
		return err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(apperror.BucketNotFound, name)
	}

	return s.repo.DropBucket(name)
}

func (s *BucketService) GetValue(name string, key string) (model.Object, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	if err := valid.Key(key); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)
	if err != nil {
		return nil, apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, apperror.New(apperror.BucketNotFound, name)
	}

	object, err := bucket.Read(key)
	if err != nil {
		return nil, apperror.WithCause(apperror.UnexpectedError, err)
	}

	if object == nil {
		return nil, apperror.New(apperror.KeyNotFound, key, name)
	}

	return object, nil
}

func (s *BucketService) PutValue(name string, key string, value model.Object) error {
	if err := valid.BucketName(name); err != nil {
		return err
	}

	if err := valid.Key(key); err != nil {
		return err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(apperror.BucketNotFound, name)
	}

	err = valid.Object(value, bucket.Schema)
	if err != nil {
		return err
	}

	return bucket.Store(key, value)
}

func (s *BucketService) DeleteValue(name string, key string) error {
	if err := valid.BucketName(name); err != nil {
		return err
	}

	if err := valid.Key(key); err != nil {
		return err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(apperror.BucketNotFound, name)
	}

	return bucket.Delete(key)
}

func (s *BucketService) Search(name string, criteria url.Values) ([]string, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return nil, apperror.WithCause(apperror.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, apperror.New(apperror.BucketNotFound, name)
	}

	if err := valid.Criteria(criteria, bucket.Schema); err != nil {
		return nil, err
	}

	normalized, err := model.Convert(criteria, bucket.Schema)
	if err != nil {
		return nil, err
	}

	return bucket.FindKeys(normalized)
}
