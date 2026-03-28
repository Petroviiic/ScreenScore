package storage

import (
	"context"
	"time"
)

func NewMockStorage() *Storage {
	return &Storage{
		StatsStorage:  &StatsMockStorage{},
		DeviceStorage: &DeviceMockStorage{},
	}
}

type StatsMockStorage struct {
	GetUsersLastFunc func(ctx context.Context, userID int64) (*UsageRecord, error)
}

func (s *StatsMockStorage) GetUsersLast(context.Context, int64, string) (*UsageRecord, error) {
	return s.GetUsersLastFunc(nil, 0)
}
func (s *StatsMockStorage) AddNewRecord(context.Context, int64, int32, string, time.Time) error {
	return nil
}
func (s *StatsMockStorage) GetGroupStats(context.Context, string, time.Time) ([]*GroupStats, error) {
	return nil, nil
}

type DeviceMockStorage struct {
}

func (d *DeviceMockStorage) Update(ctx context.Context, userId int64, deviceId string, pushToken string) error {
	return nil
}
func (d *DeviceMockStorage) GetFCMTokens(context.Context, int64) ([]string, error) {
	return []string{}, nil
}
func (d *DeviceMockStorage) RequestDeviceSync(context.Context, int) ([]string, error) {
	return []string{}, nil
}
func (d *DeviceMockStorage) UpdateLastSeen(ctx context.Context, userId int64, deviceId string) error {
	return nil
}
func (d *DeviceMockStorage) DeleteFCMToken(string) error { return nil }
