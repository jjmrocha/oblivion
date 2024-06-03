package repo

import "database/sql"

type Bucket struct {
	repo   *Repo
	Name   string  `json:"name"`
	Schema []Field `json:"schema"`
}

func (b *Bucket) Store(key string, value map[string]any) error {
	exists, err := keyExists(b.repo.db, b, key)
	if err != nil {
		return err
	}

	if exists {
		return updateValue(b.repo.db, b, key, value)
	}

	return insertValue(b.repo.db, b, key, value)
}

func (b *Bucket) Read(key string) (map[string]any, error) {
	query := buildFindByKeySql(b)
	stm, err := b.repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	row := stm.QueryRow(key)

	values := valuesForScan(b.Schema)
	err = row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	obj := buildObject(b.Schema, values)
	return obj, nil
}

func (b *Bucket) Delete(key string) error {
	query := "delete from " + b.Name + " where key = ?"
	stm, err := b.repo.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stm.Exec(key)
	return err
}

func (b *Bucket) FindKeys(criteria map[string][]any) ([]string, error) {
	query, values := buildSearchQuery(b, criteria)
	stm, err := b.repo.db.Prepare(query)
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
