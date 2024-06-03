package repo

import (
	"database/sql"
	"strings"
)

func readSchema(db *sql.DB, bucket string) ([]Field, error) {
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

	schema, err := unmarshalSchema([]byte(schemaStr))
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func buildFindByKeySql(bucket *Bucket) string {
	columns := make([]string, 0, len(bucket.Schema))
	for _, field := range bucket.Schema {
		columns = append(columns, field.Name)
	}

	columnList := strings.Join(columns, ", ")
	query := "select " + columnList + " from " + bucket.Name + " where key = ?"

	return query
}

func buildObject(schema []Field, values []any) map[string]any {
	obj := make(map[string]any)

	for i, field := range schema {
		switch field.Type {
		case StringDataType:
			holder := values[i].(*sql.NullString)
			if holder.Valid {
				obj[field.Name] = holder.String
			}
		case NumberDataType:
			holder := values[i].(*sql.NullFloat64)
			if holder.Valid {
				obj[field.Name] = holder.Float64
			}
		case BoolDataType:
			holder := values[i].(*sql.NullBool)
			if holder.Valid {
				obj[field.Name] = holder.Bool
			}
		}
	}

	return obj
}

func valuesForScan(schema []Field) []any {
	values := make([]any, len(schema))

	for i, field := range schema {
		switch field.Type {
		case StringDataType:
			var holder sql.NullString
			values[i] = &holder
		case NumberDataType:
			var holder sql.NullFloat64
			values[i] = &holder
		case BoolDataType:
			var holder sql.NullBool
			values[i] = &holder
		}
	}

	return values
}

func buildSearchQuery(bucket *Bucket, criteria map[string][]any) (string, []any) {
	where := ""
	values := make([]any, 0, len(criteria))

	for field, valueList := range criteria {
		if len(where) > 0 {
			where += " and "
		}

		or := ""

		for _, option := range valueList {
			if len(or) > 0 {
				or += " or "
			}

			or += field + " = ?"
			values = append(values, option)
		}

		where += "(" + or + ")"
	}

	query := "select key from " + bucket.Name

	if len(where) > 0 {
		query += " where " + where
	}

	return query, values
}

func bucketExists(db *sql.DB, bucket string) (bool, error) {
	schema, err := readSchema(db, bucket)
	if err != nil {
		return false, err
	}

	exists := schema != nil
	return exists, nil
}

func createTable(tx *sql.Tx, tableName string, schema []Field) error {
	query := "create table " + tableName + " (key varchar(50) primary key"
	for _, field := range schema {
		query += " , " + field.Name

		switch field.Type {
		case StringDataType:
			query += " text"
		case NumberDataType:
			query += " numeric"
		case BoolDataType:
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

func createIndex(tx *sql.Tx, tableName string, column string) error {
	indexName := "i_" + tableName + "_" + column
	query := "create index " + indexName + " on " + tableName + " (" + column + ")"

	_, err := tx.Exec(query)
	return err
}

func updateValue(db *sql.DB, bucket *Bucket, key string, obj map[string]any) error {
	columnList := ""
	values := make([]any, 0)

	for _, field := range bucket.Schema {
		if len(columnList) > 0 {
			columnList += ", "
		}

		value, found := obj[field.Name]

		if found {
			columnList += field.Name + " = ?"
			values = append(values, value)
		} else {
			columnList += field.Name + " = null"
		}
	}

	values = append(values, key)

	query := "update " + bucket.Name + " set " + columnList + " where key = ?"

	stm, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stm.Close()

	_, err = stm.Exec(values...)

	return err
}

func insertValue(db *sql.DB, bucket *Bucket, key string, obj map[string]any) error {
	columnCount := len(obj)

	columns := make([]string, 0, columnCount)
	values := make([]any, 0, columnCount+1)
	values = append(values, key)

	for field, value := range obj {
		columns = append(columns, field)
		values = append(values, value)
	}

	columnList := strings.Join(columns, ", ")
	paramList := strings.Join(strings.Split(strings.Repeat("?", columnCount), ""), ", ")
	query := "insert into " + bucket.Name + " (key, " + columnList + ") values (?, " + paramList + ")"

	stm, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stm.Close()

	_, err = stm.Exec(values...)

	return err
}

func keyExists(db *sql.DB, bucket *Bucket, key string) (bool, error) {
	query := "select count(*) from " + bucket.Name + " where key = ?"
	stm, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stm.Close()

	row := stm.QueryRow(key)

	var count int
	if err = row.Scan(&count); err != nil {
		return false, err
	}

	exists := count > 0
	return exists, nil
}
