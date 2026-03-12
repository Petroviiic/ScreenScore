package storage

import (
	"context"
	"database/sql"
)

type Storage struct {
	UserStorage interface {
		GetById(context.Context, int64) (*User, error)
	}
	StatsStorage interface {
		GetUsersLast(context.Context, int64) (*UsageRecord, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage:  &UserStorage{db},
		StatsStorage: &StatsStorage{db},
	}
}
