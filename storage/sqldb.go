package storage

import "github.com/jjmrocha/oblivion/bucket/model"

type SQLRepository interface {
}

type SQLDBRepo struct {
	impl SQLRepository
}

func NewSQLDBRepo(sql SQLRepository) *SQLDBRepo {
	repo := SQLDBRepo{
		impl: sql,
	}

	return &repo
}

func (s *SQLDBRepo) GetAllBuckets() ([]*model.Bucket, error) {
	panic("Not implemented")
}

func (s *SQLDBRepo) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	panic("Not implemented")
}

func (s *SQLDBRepo) GetBucket(name string) (*model.Bucket, error) {
	panic("Not implemented")
}

func (s *SQLDBRepo) DropBucket(name string) error {
	panic("Not implemented")
}

func (s *SQLDBRepo) Store(bucket *model.Bucket, key string, value any) error {
	panic("Not implemented")
}

func (s *SQLDBRepo) Read(bucket *model.Bucket, key string) (any, error) {
	panic("Not implemented")
}

func (s *SQLDBRepo) Delete(bucket *model.Bucket, key string) error {
	panic("Not implemented")
}

func (s *SQLDBRepo) FindKeys(bucket *model.Bucket, query map[string][]any) ([]string, error) {
	panic("Not implemented")
}
