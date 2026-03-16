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
		RegisterUser(ctx context.Context, user *User) (int64, error)
	}
	StatsStorage interface {
		GetUsersLast(context.Context, int64, string) (*UsageRecord, error)
		AddNewRecord(context.Context, int64, int32, string, time.Time) error
		GetGroupStats(context.Context, string) ([]*GroupStats, error)
	}
	GroupStorage interface {
		CheckIfMember(context.Context, int64, string) bool
		CreateGroup(context.Context, string) (string, error)
		JoinGroup(ctx context.Context, userId int64, inviteCode string) (string, error)
		LeaveGroup(ctx context.Context, userId int64, groupId string) error
		KickUser(context.Context, int64, string) error
	}
	DeviceStorage interface {
		Update(ctx context.Context, userId int64, deviceId string, pushToken string) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage:   &UserStorage{db},
		StatsStorage:  &StatsStorage{db},
		GroupStorage:  &GroupStorage{db},
		DeviceStorage: &DeviceStorage{db},
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
