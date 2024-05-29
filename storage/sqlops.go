package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	"github.com/jjmrocha/oblivion/bucket/model"
	"github.com/jjmrocha/oblivion/bucket/model/apperror"
)

func createCatalogIfNotExist(db *sql.DB) error {
	query := `create table if not exists oblivion (
				bucket_name varchar(30) primary key, 
				schema text not null
			)`

	_, err := db.Exec(query)

	return err
}

func createBucket(db *sql.DB, bucket string, schema []model.Field) error {
	exists, err := bucketExists(db, bucket)
	if err != nil {
		return err
	}

	if exists {
		return apperror.New(model.BucketAlreadyExits, bucket)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = addBucketToCatalog(tx, bucket, schema)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = createTable(tx, bucket, schema)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, field := range schema {
		if field.Indexed {
			err = createIndex(tx, bucket, field.Name)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error creating bucket %v: %v\n", bucket, err)
		return err
	}

	return nil
}

func listBuckets(db *sql.DB) ([]string, error) {
	stm, err := db.Prepare("select bucket_name from oblivion")
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

func readSchema(db *sql.DB, bucket string) ([]model.Field, error) {
	stm, err := db.Prepare("select schema from oblivion where bucket_name = ?")
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	row := stm.QueryRow(bucket)

	var schemaStr string
	if err = row.Scan(&schemaStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	schema := make([]model.Field, 0)
	err = json.Unmarshal([]byte(schemaStr), &schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func deleteBucket(db *sql.DB, bucket string) error {
	exists, err := bucketExists(db, bucket)
	if err != nil {
		return err
	}

	if !exists {
		return apperror.New(model.BucketNotFound, bucket)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = removeBucketFromCatalog(tx, bucket)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dropTable(tx, bucket)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error removing bucket %v: %v\n", bucket, err)
		return err
	}

	return nil
}

func upsertKey(db *sql.DB, bucket string, key string, value map[string]any) error {
	panic("Not done")
}

func findKey(db *sql.DB, bucket string, key string) (map[string]any, error) {
	schema, err := readSchema(db, bucket)
	if err != nil {
		return nil, err
	}

	columns := make([]string, 0)
	for _, field := range schema {
		columns = append(columns, field.Name)
	}

	columnList := strings.Join(columns, ", ")
	query := "select " + columnList + " from " + bucket + " where key = ?"

	stm, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	row := stm.QueryRow(key)

	values := make([]any, len(columns))

	err = row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	obj := make(map[string]any)

	for i, column := range columns {
		obj[column] = values[i]
	}

	return obj, nil
}

func deleteKey(db *sql.DB, bucket string, key string) error {
	query := "delete from " + bucket + " where key = ?"

	stm, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stm.Exec(key)
	return err
}

func search(db *sql.DB, bucket string, query map[string][]any) ([]string, error) {
	panic("Not done")
}

func bucketExists(db *sql.DB, bucket string) (bool, error) {
	schema, err := readSchema(db, bucket)
	if err != nil {
		return false, err
	}

	exists := schema != nil
	return exists, nil
}

func createTable(tx *sql.Tx, tableName string, schema []model.Field) error {
	query := "create table " + tableName + " (key varchar(50) primary key"
	for _, field := range schema {
		query += " , " + field.Name

		switch field.Type {
		case model.StringDataType:
			query += " text"
		case model.NumberDataType:
			query += " numeric"
		case model.BoolDataType:
			query += " boolean"
		}

		if field.Required {
			query += " not null"
		}
	}
	query += ")"

	_, err := tx.Exec(query)
	return err
}

func dropTable(tx *sql.Tx, tableName string) error {
	query := "drop table " + tableName

	_, err := tx.Exec(query)
	return err
}

func addBucketToCatalog(tx *sql.Tx, tableName string, schema []model.Field) error {
	stm, err := tx.Prepare("insert into oblivion (bucket_name, schema) values (?, ?)")
	if err != nil {
		return err
	}
	defer stm.Close()

	schemaStr, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	_, err = stm.Exec(tableName, string(schemaStr))
	return err
}

func removeBucketFromCatalog(tx *sql.Tx, tableName string) error {
	stm, err := tx.Prepare("delete from oblivion where bucket_name = ?")
	if err != nil {
		return err
	}
	defer stm.Close()

	_, err = stm.Exec(tableName)
	return err
}

func createIndex(tx *sql.Tx, tableName string, column string) error {
	indexName := "i_" + tableName + "_" + column
	query := "create index " + indexName + " on " + tableName + " (" + column + ")"

	_, err := tx.Exec(query)
	return err
}
