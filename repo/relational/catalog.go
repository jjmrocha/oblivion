package relational

import (
	"database/sql"

	"github.com/jjmrocha/oblivion/model"
)

func createCatalogIfNotExist(db *sql.DB) error {
	query := `create table if not exists oblivion (
				bucket_name varchar(30) primary key, 
				schema text not null
			)`

	_, err := db.Exec(query)

	return err
}

func addBucketToCatalog(tx *sql.Tx, bucket string, schema []model.Field) error {
	stm, err := tx.Prepare("insert into oblivion (bucket_name, schema) values (?, ?)")
	if err != nil {
		return err
	}
	defer stm.Close()

	data, err := marshalSchema(schema)
	if err != nil {
		return err
	}

	_, err = stm.Exec(bucket, string(data))
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

func bucketList(db *sql.DB) ([]string, error) {
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
