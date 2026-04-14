package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StatsStorage struct {
	db *sql.DB
}

type UsageRecord struct {
	ScreenTime int32     `json:"screen_time"`
	DeviceID   string    `json:"device_id"`
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type GroupStats struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	ScreenTime int32     `json:"screen_time"`
	RecordedAt time.Time `json:"recorded_at"`
}

func (s *StatsStorage) GetUsersLast(ctx context.Context, userId int64, deviceID string) (*UsageRecord, error) {
	query := `	SELECT screen_time, device_id, recorded_at, created_at FROM screen_time_logs 
				WHERE user_id = $1 AND device_id = $2
				ORDER BY recorded_at DESC LIMIT 1;`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	record := &UsageRecord{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
		deviceID,
	).Scan(
		&record.ScreenTime,
		&record.DeviceID,
		&record.RecordedAt,
		&record.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *StatsStorage) AddNewRecord(ctx context.Context, userId int64, screenTime int32, deviceId string, recordedAt time.Time) error {
	query := ` 	INSERT INTO screen_time_logs(user_id, screen_time, device_id, recorded_at) 
				VALUES($1, $2, $3, $4);
			`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		userId,
		screenTime,
		deviceId,
		recordedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *StatsStorage) GetGroupStats(ctx context.Context, groupId string, desiredDate time.Time) ([]*GroupStats, error) {
	query := `
			WITH group_users AS (
				SELECT user_id FROM group_members WHERE group_id = $1
			),
			ranked_stats AS (
				SELECT user_id, screen_time, recorded_at,
					ROW_NUMBER() OVER (PARTITION BY device_id,user_id ORDER BY recorded_at DESC, screen_time DESC) as rn
				FROM screen_time_logs
				WHERE user_id IN (SELECT user_id FROM group_users)
				AND recorded_at::DATE BETWEEN $2 AND $3
			)
			SELECT 
				gu.user_id, 
				u.email, 
				u.username, 
				COALESCE(SUM(rs.screen_time), 0) as total_screen_time,
				MAX(rs.recorded_at) as last_recorded_at
			FROM group_users gu
			JOIN users u ON u.id = gu.user_id
			JOIN ranked_stats rs ON rs.user_id = gu.user_id AND rs.rn = 1 
			GROUP BY gu.user_id, u.email, u.username
			ORDER BY total_screen_time DESC; 
	`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	startTime := desiredDate.UTC()
	endTime := startTime.AddDate(0, 0, 1)

	rows, err := s.db.QueryContext(
		ctx,
		query,
		groupId,
		startTime,
		endTime,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*GroupStats
	for rows.Next() {
		stat := &GroupStats{}
		err := rows.Scan(
			&stat.ID,
			&stat.Email,
			&stat.Username,
			&stat.ScreenTime,
			&stat.RecordedAt,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

func (s *StatsStorage) GetUserAverageScreenTimeForWeek(ctx context.Context, weekStart time.Time, weekEnd time.Time, userID int64) (float64, error) {
	query := `
		WITH ranked_stats AS (
			SELECT user_id, screen_time, recorded_at,
				ROW_NUMBER() OVER (PARTITION BY device_id, recorded_at::DATE ORDER BY recorded_at DESC, screen_time DESC) as rn
			FROM screen_time_logs
			WHERE user_id = 16	
			AND recorded_at >= $1::TIMESTAMP AND recorded_at < $2
		)
		SELECT 
			COALESCE(AVG(rs.screen_time), 0) as average_screen_time
		FROM ranked_stats rs
		WHERE rs.rn = 1;
	`

	averageMins := -1.0
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
		weekStart,
		weekEnd,
	).Scan(
		&averageMins,
	)

	if err != nil {
		return -1, err
	}
	if averageMins == -1 {
		return -1, fmt.Errorf("unknown error")
	}

	return averageMins, nil
}
func (s *StatsStorage) GetUserScreenTimeForDay(ctx context.Context, date time.Time, userID int64) (int, error) {
	query := `
		WITH ranked_stats AS (
			SELECT user_id, screen_time, recorded_at,
				ROW_NUMBER() OVER (PARTITION BY device_id,user_id ORDER BY recorded_at DESC, screen_time DESC) as rn
			FROM screen_time_logs
			WHERE user_id = $1
			AND recorded_at >= $2::TIMESTAMP AND recorded_at < $2::TIMESTAMP + interval '1 day'
		)
		SELECT 
			COALESCE(SUM(rs.screen_time), 0) as total_screen_time
		FROM ranked_stats rs
		WHERE rs.rn = 1;
	`

	totalMins := -1
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
		date,
	).Scan(
		&totalMins,
	)

	if err != nil {
		return -1, err
	}
	if totalMins == -1 {
		return -1, fmt.Errorf("unknown error")
	}

	return totalMins, nil
}
