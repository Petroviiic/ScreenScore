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
}

func (s *StatsMockStorage) GetUsersLast(context.Context, int64) (*UsageRecord, error) {
	return &UsageRecord{}, nil
}
func (s *StatsMockStorage) AddNewRecord(context.Context, int64, int32, time.Time) error {
	return nil
}
