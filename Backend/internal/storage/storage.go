package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ERROR_NO_ROWS_AFFECTED    = errors.New("no rows affected")
	ERROR_DUPLICATE_KEY_VALUE = errors.New("record already exists")
)

type Storage struct {
	UserStorage interface {
		GetById(context.Context, int64) (*User, error)
		GetByUsername(ctx context.Context, username string) (*User, error)
		RegisterUser(ctx context.Context, user *User) (int64, error)
	}
	StatsStorage interface {
		GetUsersLast(context.Context, int64, string) (*UsageRecord, error)
		AddNewRecord(context.Context, int64, int32, string, time.Time) error
		GetGroupStats(context.Context, string, time.Time) ([]*GroupStats, error)
	}
	GroupStorage interface {
		CheckIfMember(context.Context, int64, string) bool
		CreateGroup(context.Context, string) (string, error)
		GetGroupMembersExclusive(context.Context, string, int64) ([]int, error)
		JoinGroup(ctx context.Context, userId int64, inviteCode string) (string, error)
		LeaveGroup(ctx context.Context, userId int64, groupId string) error
		KickUser(context.Context, int64, string) error
	}
	DeviceStorage interface {
		Update(ctx context.Context, userId int64, deviceId string, pushToken string) error
		GetFCMTokens(context.Context, int64) ([]string, error)
		RequestDeviceSync(context.Context, int) ([]string, error)
		UpdateLastSeen(ctx context.Context, userId int64, deviceId string) error
		DeleteFCMToken(string) error
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
