package storage

import (
	"context"
	"database/sql"
	"time"
)

const (
	ERROR_NO_ROWS_AFFECTED = "no rows affected"
)

type Storage struct {
	UserStorage interface {
		GetById(context.Context, int64) (*User, error)
		RegisterUser(ctx context.Context, user *User) error
	}
	StatsStorage interface {
		GetUsersLast(context.Context, int64) (*UsageRecord, error)
		AddNewRecord(context.Context, int64, int32, time.Time) error
		GetGroupStats(context.Context, string) ([]*GroupStats, error)
	}
	GroupStorage interface {
		CheckIfMember(context.Context, int64, string) bool
		CreateGroup(context.Context, string) (string, error)
		JoinGroup(ctx context.Context, userId int64, inviteCode string) error
		LeaveGroup(ctx context.Context, userId int64, groupId string) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage:  &UserStorage{db},
		StatsStorage: &StatsStorage{db},
		GroupStorage: &GroupStorage{db},
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
