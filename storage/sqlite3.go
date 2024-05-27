package storage

type SqliteRepository struct {
}

func NewSqliteRepository() *SqliteRepository {
	repo := SqliteRepository{}
	return &repo
}
