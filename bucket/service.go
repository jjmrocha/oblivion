package bucket

import (
	"strconv"

	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/repo"
	"github.com/jjmrocha/oblivion/valid"
)

type BucketService struct {
	repo *repo.Repo
}

func NewBucketService(repo *repo.Repo) *BucketService {
	service := BucketService{
		repo: repo,
	}
	return &service
}

func (s *BucketService) BucketList() ([]string, error) {
	bucketList, err := s.repo.GetAllBuckets()
	if err != nil {
		return nil, model.ErrorWithReason(model.UnexpectedError, err)
	}

	return bucketList, nil
}

func (s *BucketService) CreateBucket(name string, schema []repo.Field) (*repo.Bucket, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	if err := valid.Schema(schema); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)
	if err != nil {
		return nil, model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket != nil {
		return nil, model.Error(model.BucketAlreadyExits, name)
	}

	return s.repo.CreateBucket(name, schema)
}

func (s *BucketService) GetBucket(name string) (*repo.Bucket, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return nil, model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, model.Error(model.BucketNotFound, name)
	}

	return bucket, nil
}

func (s *BucketService) DeleteBucket(name string) error {
	if err := valid.BucketName(name); err != nil {
		return err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return model.Error(model.BucketNotFound, name)
	}

	return s.repo.DropBucket(name)
}

func (s *BucketService) GetValue(name string, key string) (repo.Object, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	if err := valid.Key(key); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)
	if err != nil {
		return nil, model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, model.Error(model.BucketNotFound, name)
	}

	object, err := bucket.Read(key)
	if err != nil {
		return nil, model.ErrorWithReason(model.UnexpectedError, err)
	}

	if object == nil {
		return nil, model.Error(model.KeyNotFound, key, name)
	}

	return object, nil
}

func (s *BucketService) PutValue(name string, key string, value repo.Object) error {
	if err := valid.BucketName(name); err != nil {
		return err
	}

	if err := valid.Key(key); err != nil {
		return err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return model.Error(model.BucketNotFound, name)
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
		return model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return model.Error(model.BucketNotFound, name)
	}

	return bucket.Delete(key)
}

func (s *BucketService) Search(name string, query map[string][]string) ([]string, error) {
	if err := valid.BucketName(name); err != nil {
		return nil, err
	}

	bucket, err := s.repo.GetBucket(name)

	if err != nil {
		return nil, model.ErrorWithReason(model.UnexpectedError, err)
	}

	if bucket == nil {
		return nil, model.Error(model.BucketNotFound, name)
	}

	normalized, err := normalize(query, bucket.Schema)
	if err != nil {
		return nil, err
	}

	return bucket.FindKeys(normalized)
}

func normalize(query map[string][]string, schema []repo.Field) (map[string][]any, error) {
	fieldMap := make(map[string]repo.Field)
	for _, field := range schema {
		fieldMap[field.Name] = field
	}

	normalized := make(map[string][]any)

	for name, values := range query {
		field, found := fieldMap[name]
		if !found {
			return nil, model.Error(model.UnknownField, name)
		}

		switch field.Type {
		case repo.StringDataType:
			normalized[name] = convertStrings(values)
		case repo.NumberDataType:
			floats, err := convertFloats(values)
			if err != nil {
				return nil, model.Error(model.InvalidField, name)
			}

			normalized[name] = floats
		case repo.BoolDataType:
			bools, err := convertBools(values)
			if err != nil {
				return nil, model.Error(model.InvalidField, name)
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
