package relational

import (
	"database/sql"
	"strings"

	"github.com/jjmrocha/oblivion/model"
)

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

	schema, err := unmarshalSchema([]byte(schemaStr))
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func buildFindByKeySql(bucket *bucket) string {
	columns := make([]string, 0, len(bucket.schema))
	for _, field := range bucket.schema {
		columns = append(columns, field.Name)
	}

	columnList := strings.Join(columns, ", ")
	query := "select " + columnList + " from " + bucket.name + " where key = ?"

	return query
}

func buildObject(schema []model.Field, values []any) model.Object {
	obj := make(model.Object)

	for i, field := range schema {
		switch field.Type {
		case model.StringDataType:
			holder := values[i].(*sql.NullString)
			if holder.Valid {
				obj[field.Name] = holder.String
			}
		case model.NumberDataType:
			holder := values[i].(*sql.NullFloat64)
			if holder.Valid {
				obj[field.Name] = holder.Float64
			}
		case model.BoolDataType:
			holder := values[i].(*sql.NullBool)
			if holder.Valid {
				obj[field.Name] = holder.Bool
			}
		}
	}

	return obj
}

func valuesForScan(schema []model.Field) []any {
	values := make([]any, len(schema))

	for i, field := range schema {
		switch field.Type {
		case model.StringDataType:
			var holder sql.NullString
			values[i] = &holder
		case model.NumberDataType:
			var holder sql.NullFloat64
			values[i] = &holder
		case model.BoolDataType:
			var holder sql.NullBool
			values[i] = &holder
		}
	}

	return values
}

func buildSearchQuery(bucket *bucket, criteria model.Criteria) (string, []any) {
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

	query := "select key from " + bucket.name

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

func createIndex(tx *sql.Tx, tableName string, column string) error {
	indexName := "i_" + tableName + "_" + column
	query := "create index " + indexName + " on " + tableName + " (" + column + ")"

	_, err := tx.Exec(query)
	return err
}

func updateValue(db *sql.DB, bucket *bucket, key string, obj model.Object) error {
	columnList := ""
	values := make([]any, 0)

	for _, field := range bucket.schema {
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

	query := "update " + bucket.name + " set " + columnList + " where key = ?"

	stm, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stm.Close()

	_, err = stm.Exec(values...)

	return err
}

func insertValue(db *sql.DB, bucket *bucket, key string, obj model.Object) error {
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
	query := "insert into " + bucket.name + " (key, " + columnList + ") values (?, " + paramList + ")"

	stm, err := db.Prepare(query)
	if err != nil {
		return err
	}

	defer stm.Close()

	_, err = stm.Exec(values...)

	return err
}

func keyExists(db *sql.DB, bucket *bucket, key string) (bool, error) {
	query := "select count(*) from " + bucket.name + " where key = ?"
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
