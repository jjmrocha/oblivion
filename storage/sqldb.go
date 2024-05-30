package storage

import (
	"database/sql"
	"log"

	"github.com/jjmrocha/oblivion/bucket/model"
)

type SQLDBRepo struct {
	db *sql.DB
}

func NewSQLDBRepo(driver string, datasource string) *SQLDBRepo {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		log.Panicf("Error opening db %v using driver %v: %v", datasource, driver, err)
	}

	err = createCatalogIfNotExist(db)
	if err != nil {
		log.Panicf("Error creating db catalog on %v using driver %v: %v", datasource, driver, err)
	}

	repo := SQLDBRepo{
		db: db,
	}

	return &repo
}

func (r *SQLDBRepo) Close() {
	if err := r.db.Close(); err != nil {
		log.Printf("Error closing db: %v\n", err)
	}
}

func (r *SQLDBRepo) GetAllBuckets() ([]string, error) {
	return listBuckets(r.db)
}

func (r *SQLDBRepo) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	if err := createBucket(r.db, name, schema); err != nil {
		return nil, err
	}

	bucket := model.Bucket{
		Name:   name,
		Schema: schema,
	}

	return &bucket, nil
}

func (r *SQLDBRepo) GetBucket(name string) (*model.Bucket, error) {
	schema, err := readSchema(r.db, name)
	if err != nil {
		return nil, err
	}

	if schema == nil {
		return nil, nil
	}

	bucket := model.Bucket{
		Name:   name,
		Schema: schema,
	}

	return &bucket, nil
}

func (r *SQLDBRepo) DropBucket(name string) error {
	return deleteBucket(r.db, name)
}

func (r *SQLDBRepo) Store(bucket *model.Bucket, key string, value map[string]any) error {
	return upsertKey(r.db, bucket, key, value)
}

func (r *SQLDBRepo) Read(bucket *model.Bucket, key string) (map[string]any, error) {
	return findKey(r.db, bucket, key)
}

func (r *SQLDBRepo) Delete(bucket *model.Bucket, key string) error {
	return deleteKey(r.db, bucket, key)
}

func (r *SQLDBRepo) FindKeys(bucket *model.Bucket, query map[string][]any) ([]string, error) {
	return search(r.db, bucket, query)
}
