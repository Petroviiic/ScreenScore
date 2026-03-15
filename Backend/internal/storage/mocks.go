package storage

import (
	"context"
	"time"
)

func NewMockStorage() *Storage {
	return &Storage{
		StatsStorage: &StatsMockStorage{},
	}
}

type StatsMockStorage struct {
	GetUsersLastFunc func(ctx context.Context, userID int64) (*UsageRecord, error)
}

func (s *StatsMockStorage) GetUsersLast(context.Context, int64) (*UsageRecord, error) {
	return s.GetUsersLastFunc(nil, 0)
}
func (s *StatsMockStorage) AddNewRecord(context.Context, int64, int32, time.Time) error {
	return nil
}
func (s *StatsMockStorage) GetGroupStats(context.Context, string) ([]*GroupStats, error) {
	return nil, nil
}
