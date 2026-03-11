package storage

import "database/sql"

type Storage struct {
	UserStorage interface {
		GetById()
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage: &UserStorage{db},
	}
}
