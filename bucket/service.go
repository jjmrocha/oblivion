package bucket

import (
	"context"
	"net/url"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
	"github.com/jjmrocha/oblivion/valid"
)

type BucketService struct {
	repo repo.Repository
}

func NewService(repo repo.Repository) *BucketService {
	service := BucketService{
		repo: repo,
	}
	return &service
}

func (s *BucketService) BucketList(ctx context.Context) ([]string, error) {
	bucketList, err := s.repo.BucketNames(ctx)
	if err != nil {
		return nil, apperror.UnexpectedError.WithCause(err)
	}

	return bucketList, nil
}

func (s *BucketService) CreateBucket(ctx context.Context, name string, schema []model.Field) (repo.Bucket, error) {
	bucket, err := s.repo.GetBucket(ctx, name)
	if err != nil {
		return nil, apperror.UnexpectedError.WithCause(err)
	}

	if bucket != nil {
		return nil, apperror.BucketAlreadyExits.New(name)
	}

	return s.repo.NewBucket(ctx, name, schema)
}

func (s *BucketService) GetBucket(ctx context.Context, name string) (repo.Bucket, error) {
	bucket, err := s.repo.GetBucket(ctx, name)

	if err != nil {
		return nil, apperror.UnexpectedError.WithCause(err)
	}

	if bucket == nil {
		return nil, apperror.BucketNotFound.New(name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(ctx context.Context, name string) error {
	bucket, err := s.repo.GetBucket(ctx, name)

	if err != nil {
		return apperror.UnexpectedError.WithCause(err)
	}

	if bucket == nil {
		return apperror.BucketNotFound.New(name)
	}

	return s.repo.DropBucket(ctx, name)
}

func (s *BucketService) Value(ctx context.Context, name string, key string) (model.Object, error) {
	bucket, err := s.repo.GetBucket(ctx, name)
	if err != nil {
		return nil, apperror.UnexpectedError.WithCause(err)
	}

	if bucket == nil {
		return nil, apperror.BucketNotFound.New(name)
	}

	object, err := bucket.Read(ctx, key)
	if err != nil {
		return nil, apperror.UnexpectedError.WithCause(err)
	}

	if object == nil {
		return nil, apperror.KeyNotFound.New(key, name)
	}

	return object, nil
}

func (s *BucketService) SetValue(ctx context.Context, name string, key string, value model.Object) error {
	bucket, err := s.repo.GetBucket(ctx, name)

	if err != nil {
		return apperror.UnexpectedError.WithCause(err)
	}

	if bucket == nil {
		return apperror.BucketNotFound.New(name)
	}

	err = valid.Object(value, bucket.Schema())
	if err != nil {
		return err
	}

	return bucket.Store(ctx, key, value)
}

func (s *BucketService) DeleteValue(ctx context.Context, name string, key string) error {
	bucket, err := s.repo.GetBucket(ctx, name)

	if err != nil {
		return apperror.UnexpectedError.WithCause(err)
	}

	if bucket == nil {
		return apperror.BucketNotFound.New(name)
	}

	return bucket.Delete(ctx, key)
}

func (s *BucketService) FindKeys(ctx context.Context, name string, criteria url.Values) ([]string, error) {
	bucket, err := s.repo.GetBucket(ctx, name)

	if err != nil {
		return nil, apperror.UnexpectedError.WithCause(err)
	}

	if bucket == nil {
		return nil, apperror.BucketNotFound.New(name)
	}

	if err := valid.Criteria(criteria, bucket.Schema()); err != nil {
		return nil, err
	}

	normalized, err := model.Convert(criteria, bucket.Schema())
	if err != nil {
		return nil, err
	}

	return bucket.Keys(ctx, normalized)
}
