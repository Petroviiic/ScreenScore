package storage

import (
	"context"
	"database/sql"
	"time"
)

type Storage struct {
	UserStorage interface {
		GetById(context.Context, int64) (*User, error)
	}
	StatsStorage interface {
		GetUsersLast(context.Context, int64) (*UsageRecord, error)
		AddNewRecord(context.Context, int64, int32, time.Time) error
		GetGroupStats(context.Context, string) ([]*GroupStats, error)
	}
	GroupsStorage interface {
		CheckIfMember(context.Context, int64, string) bool
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage:   &UserStorage{db},
		StatsStorage:  &StatsStorage{db},
		GroupsStorage: &GroupsStorage{db},
	}
}
