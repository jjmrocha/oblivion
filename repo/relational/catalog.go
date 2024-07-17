package relational

import (
	"context"
	"database/sql"

	"github.com/jjmrocha/oblivion/model"
)

func createCatalogIfNotExist(ctx context.Context, db *sql.DB) error {
	query := `create table if not exists oblivion (
				bucket_name varchar(30) primary key, 
				schema text not null
			)`

	_, err := db.ExecContext(ctx, query)

	return err
}

func addBucketToCatalog(ctx context.Context, tx *sql.Tx, bucket string, schema []model.Field) error {
	stm, err := tx.PrepareContext(ctx, "insert into oblivion (bucket_name, schema) values (?, ?)")
	if err != nil {
		return err
	}
	defer stm.Close()

	data, err := marshalSchema(schema)
	if err != nil {
		return err
	}

	_, err = stm.ExecContext(ctx, bucket, string(data))
	return err
}

func removeBucketFromCatalog(ctx context.Context, tx *sql.Tx, tableName string) error {
	stm, err := tx.PrepareContext(ctx, "delete from oblivion where bucket_name = ?")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.ExecContext(ctx, tableName)
	return err
}

func bucketList(ctx context.Context, db *sql.DB) ([]string, error) {
	stm, err := db.PrepareContext(ctx, "select bucket_name from oblivion")
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.QueryContext(ctx)
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
