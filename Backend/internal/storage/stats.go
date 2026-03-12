package storage

import (
	"context"
	"database/sql"
	"time"
)

type StatsStorage struct {
	db *sql.DB
}

type UsageRecord struct {
	ScreenTime int32     `json:"screen_time"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *StatsStorage) GetUsersLast(ctx context.Context, userId int64) (*UsageRecord, error) {
	query := `	SELECT screen_time, recorded_at, created_at FROM screen_time_logs 
				WHERE user_id = $1
				ORDER BY recorded_at DESC LIMIT 1;`

	record := &UsageRecord{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&record.ScreenTime,
		&record.RecordedAt,
		&record.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *StatsStorage) AddNewRecord(ctx context.Context, userId int64, screenTime int32, recordedAt time.Time) error {
	query := ` 	INSERT INTO screen_time_logs(user_id, screen_time, recorded_at) 
				VALUES($1, $2, $3);
			`

	_, err := s.db.ExecContext(
		ctx,
		query,
		userId,
		screenTime,
		recordedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
