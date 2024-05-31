package storage

import (
	"database/sql"
	"log"

	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
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
	stm, err := r.db.Prepare("select bucket_name from oblivion")
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bucketList := make([]string, 0)
	var bucket string

	for rows.Next() {
		if err = rows.Scan(&bucket); err != nil {
			return nil, err
		}

		bucketList = append(bucketList, bucket)
	}

	return bucketList, nil
}

func (r *SQLDBRepo) CreateBucket(name string, schema []model.Field) (*model.Bucket, error) {
	exists, err := bucketExists(r.db, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, apperror.New(model.BucketAlreadyExits, name)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	err = addBucketToCatalog(tx, name, schema)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = createTable(tx, name, schema)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, field := range schema {
		if field.Indexed {
			err = createIndex(tx, name, field.Name)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error creating bucket %v: %v\n", name, err)
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
	exists, err := bucketExists(r.db, name)
	if err != nil {
		return err
	}

	if !exists {
		return apperror.New(model.BucketNotFound, name)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = removeBucketFromCatalog(tx, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dropTable(tx, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error removing bucket %v: %v\n", name, err)
		return err
	}

	return nil
}

func (r *SQLDBRepo) Store(bucket *model.Bucket, key string, value map[string]any) error {
	exists, err := keyExists(r.db, bucket, key)
	if err != nil {
		return err
	}

	if exists {
		return updateValue(r.db, bucket, key, value)
	}

	return insertValue(r.db, bucket, key, value)
}

func (r *SQLDBRepo) Read(bucket *model.Bucket, key string) (map[string]any, error) {
	query := buildFindByKeySql(bucket)
	stm, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	row := stm.QueryRow(key)

	values := valuesForScan(bucket)
	err = row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	obj := buildObject(bucket, values)
	return obj, nil
}

func (r *SQLDBRepo) Delete(bucket *model.Bucket, key string) error {
	query := "delete from " + bucket.Name + " where key = ?"
	stm, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stm.Exec(key)
	return err
}

func (r *SQLDBRepo) FindKeys(bucket *model.Bucket, criteria map[string][]any) ([]string, error) {
	query, values := buildSearchQuery(bucket, criteria)
	stm, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.Query(values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keyList := make([]string, 0)
	var key string

	for rows.Next() {
		if err = rows.Scan(&key); err != nil {
			return nil, err
		}

		keyList = append(keyList, key)
	}

	return keyList, nil
}
