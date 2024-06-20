package repo

import (
	"database/sql"
	"log"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/model"
)

type Repo struct {
	db *sql.DB
}

func New(driver string, datasource string) *Repo {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		log.Panicf("Error opening db %v using driver %v: %v", datasource, driver, err)
	}

	err = createCatalogIfNotExist(db)
	if err != nil {
		log.Panicf("Error creating db catalog on %v using driver %v: %v", datasource, driver, err)
	}

	repo := Repo{
		db: db,
	}

	return &repo
}

func (r *Repo) Close() {
	if err := r.db.Close(); err != nil {
		log.Printf("Error closing db: %v\n", err)
	}
}

func (r *Repo) BucketNames() ([]string, error) {
	return bucketList(r.db)
}

func (r *Repo) NewBucket(name string, schema []model.Field) (*Bucket, error) {
	exists, err := bucketExists(r.db, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, apperror.BucketAlreadyExits.NewError(name)
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

	bucket := Bucket{
		repo:   r,
		Name:   name,
		Schema: schema,
	}

	return &bucket, nil
}

func (r *Repo) GetBucket(name string) (*Bucket, error) {
	schema, err := readSchema(r.db, name)
	if err != nil {
		return nil, err
	}

	if schema == nil {
		return nil, nil
	}

	bucket := Bucket{
		repo:   r,
		Name:   name,
		Schema: schema,
	}

	return &bucket, nil
}

func (r *Repo) DropBucket(name string) error {
	exists, err := bucketExists(r.db, name)
	if err != nil {
		return err
	}

	if !exists {
		return apperror.BucketNotFound.NewError(name)
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
