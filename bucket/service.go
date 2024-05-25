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

func (s *BucketService) PutValue(name string, key string, value map[string]any) error {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return apperror.New(model.BucketNotFound, name)
	}

	err = checkValue(value, bucket.Schema)
	if err != nil {
		return err
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

func checkValue(object map[string]any, schema []model.Field) error {
	fieldMap := make(map[string]model.Field)
	for _, field := range schema {
		fieldMap[field.Name] = field
	}

	for name, value := range object {
		field, found := fieldMap[name]
		if !found {
			return apperror.New(model.UnknownField, name)
		}

		switch field.Type {
		case model.StringDataType:
			if _, ok := value.(string); !ok {
				return apperror.New(model.InvalidField, name)
			}
		case model.NumberDataType:
			if _, ok := value.(float64); !ok {
				return apperror.New(model.InvalidField, name)
			}
		case model.BoolDataType:
			if _, ok := value.(bool); !ok {
				return apperror.New(model.InvalidField, name)
			}
		}

		for _, field := range schema {
			if !field.Required {
				continue
			}

			if _, found := object[field.Name]; !found {
				return apperror.New(model.MissingField, field.Name)
			}
		}
	}

	return nil
}
