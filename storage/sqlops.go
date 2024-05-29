package storage

import (
	"database/sql"

	"github.com/jjmrocha/oblivion/bucket/model"
)

func createCatalogIfNotExist(db *sql.DB) error {
	sql := `create table if not exists oblivion (
							bucket_name text primary key, 
							schema text not null
					)`

	_, err := db.Exec(sql)

	return err
}

func createTable(db *sql.DB, bucket string, schema []model.Field) error {
	panic("Not done")
}

func listBuckets(db *sql.DB) ([]string, error) {
	panic("Not done")
}

func readSchema(db *sql.DB, bucket string) ([]model.Field, error) {
	panic("Not done")
}

func deleteTable(db *sql.DB, bucket string) error {
	panic("Not done")
}

func upsertKey(db *sql.DB, bucket string, key string, value map[string]any) error {
	panic("Not done")
}

func findKey(db *sql.DB, bucket string, key string) (map[string]any, error) {
	panic("Not done")
}

func deleteKey(db *sql.DB, bucket string, key string) error {
	panic("Not done")
}

func search(db *sql.DB, bucket string, query map[string][]any) ([]string, error) {
	panic("Not done")
}
