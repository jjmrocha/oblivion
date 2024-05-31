package bucket

import (
	"regexp"
	"strconv"

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

	return bucketList, nil
}

func (s *BucketService) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	matched, err := regexp.MatchString("^[a-zA-Z][a-zA-Z_0-9]*[a-zA-Z0-9]$", name)
	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	if !matched || len(name) > 30 {
		return nil, apperror.New(model.InvalidBucketName, name)
	}

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

	object, err := s.repository.Read(bucket, key)
	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	if object == nil {
		return nil, apperror.New(model.KeyNotFound, key, name)
	}

	return object, nil
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

func (s *BucketService) Search(name string, query map[string][]string) ([]string, error) {
	bucket, err := s.repository.GetBucket(name)

	if err != nil {
		return nil, apperror.WithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, apperror.New(model.BucketNotFound, name)
	}

	normalized, err := normalize(query, bucket.Schema)
	if err != nil {
		return nil, err
	}

	return s.repository.FindKeys(bucket, normalized)
}

func normalize(query map[string][]string, schema []model.Field) (map[string][]any, error) {
	fieldMap := make(map[string]model.Field)
	for _, field := range schema {
		fieldMap[field.Name] = field
	}

	normalized := make(map[string][]any)

	for name, values := range query {
		field, found := fieldMap[name]
		if !found {
			return nil, apperror.New(model.UnknownField, name)
		}

		switch field.Type {
		case model.StringDataType:
			normalized[name] = convertStrings(values)
		case model.NumberDataType:
			floats, err := convertFloats(values)
			if err != nil {
				return nil, apperror.New(model.InvalidField, name)
			}

			normalized[name] = floats
		case model.BoolDataType:
			bools, err := convertBools(values)
			if err != nil {
				return nil, apperror.New(model.InvalidField, name)
			}

			normalized[name] = bools
		}
	}

	return normalized, nil
}

func convertFloats(input []string) ([]any, error) {
	values := make([]any, 0)

	for _, strValue := range input {
		value, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
}

func convertStrings(input []string) []any {
	values := make([]any, 0)

	for _, value := range input {
		values = append(values, value)
	}

	return values
}

func convertBools(input []string) ([]any, error) {
	values := make([]any, 0)

	for _, strValue := range input {
		value, err := strconv.ParseBool(strValue)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	return values, nil
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
