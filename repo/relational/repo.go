package relational

import (
	"context"
	"database/sql"
	"log"

	"github.com/jjmrocha/oblivion/apperror"
	"github.com/jjmrocha/oblivion/model"
	"github.com/jjmrocha/oblivion/repo"
)

type sqlRepo struct {
	db *sql.DB
}

func New(driver string, datasource string) repo.Repository {
	db, err := sql.Open(driver, datasource)
	if err != nil {
		log.Panicf("Error opening db %v using driver %v: %v", datasource, driver, err)
	}

	err = createCatalogIfNotExist(context.Background(), db)
	if err != nil {
		log.Panicf("Error creating db catalog on %v using driver %v: %v", datasource, driver, err)
	}

	repo := sqlRepo{
		db: db,
	}

	return &repo
}

func (r *sqlRepo) Close() {
	if err := r.db.Close(); err != nil {
		log.Printf("Error closing db: %v\n", err)
	}
}

func (r *sqlRepo) BucketNames(ctx context.Context) ([]string, error) {
	return bucketList(ctx, r.db)
}

func (r *sqlRepo) NewBucket(ctx context.Context, name string, schema []model.Field) (repo.Bucket, error) {
	exists, err := bucketExists(ctx, r.db, name)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, apperror.BucketAlreadyExits.NewError(name)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = addBucketToCatalog(ctx, tx, name, schema)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = createTable(ctx, tx, name, schema)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, field := range schema {
		if field.Indexed {
			err = createIndex(ctx, tx, name, field.Name)
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

	bucket := bucket{
		repo:   r,
		name:   name,
		schema: schema,
	}

	return &bucket, nil
}

func (r *sqlRepo) GetBucket(ctx context.Context, name string) (repo.Bucket, error) {
	schema, err := readSchema(ctx, r.db, name)
	if err != nil {
		return nil, err
	}

	if schema == nil {
		return nil, nil
	}

	bucket := bucket{
		repo:   r,
		name:   name,
		schema: schema,
	}

	return &bucket, nil
}

func (r *sqlRepo) DropBucket(ctx context.Context, name string) error {
	exists, err := bucketExists(ctx, r.db, name)
	if err != nil {
		return err
	}

	if !exists {
		return apperror.BucketNotFound.NewError(name)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = removeBucketFromCatalog(ctx, tx, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dropTable(ctx, tx, name)
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
