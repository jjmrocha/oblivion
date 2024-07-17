package relational

import (
	"context"
	"database/sql"

	"github.com/jjmrocha/oblivion/model"
)

type bucket struct {
	repo   *sqlRepo
	name   string
	schema []model.Field
}

func (b *bucket) Name() string {
	return b.name
}

func (b *bucket) Schema() []model.Field {
	return b.schema
}

func (b *bucket) Store(ctx context.Context, key string, value model.Object) error {
	exists, err := keyExists(ctx, b.repo.db, b, key)
	if err != nil {
		return err
	}

	if exists {
		return updateValue(ctx, b.repo.db, b, key, value)
	}

	return insertValue(ctx, b.repo.db, b, key, value)
}

func (b *bucket) Read(ctx context.Context, key string) (model.Object, error) {
	query := buildFindByKeySql(b)
	stm, err := b.repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stm.Close()

	row := stm.QueryRowContext(ctx, key)

	values := valuesForScan(b.schema)
	err = row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	obj := buildObject(b.schema, values)
	return obj, nil
}

func (b *bucket) Delete(ctx context.Context, key string) error {
	query := "delete from " + b.name + " where key = ?"
	stm, err := b.repo.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stm.ExecContext(ctx, key)
	return err
}

func (b *bucket) Keys(ctx context.Context, criteria model.Criteria) ([]string, error) {
	query, values := buildSearchQuery(b, criteria)
	stm, err := b.repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	rows, err := stm.QueryContext(ctx, values...)
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
