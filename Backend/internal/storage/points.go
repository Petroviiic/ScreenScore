package storage

import (
	"context"
	"database/sql"
	"time"
)

type PointsLogicsStorage struct {
	db *sql.DB
}
type WeeklyGroupStats struct {
	GroupID                string
	UserID                 int64
	ScreenTime             int
	GroupAverageScreenTime float64
}

func (p *PointsLogicsStorage) GetWeeklyGroupStats(ctx context.Context, startDate, endDate time.Time) ([]*WeeklyGroupStats, error) {
	query := `
		WITH group_users AS (
			SELECT user_id, group_id FROM group_members
		),
		ranked_stats AS (
			SELECT user_id, screen_time, recorded_at,
				ROW_NUMBER() OVER (PARTITION BY device_id, user_id ORDER BY recorded_at DESC, screen_time DESC) as rn
			FROM screen_time_logs
			WHERE user_id IN (SELECT user_id FROM group_users)
			AND recorded_at::DATE between $1 and $2
		),
		user_totals AS (
			SELECT 
				gu.group_id,
				gu.user_id, 
				u.email, 
				u.username, 
				COALESCE(SUM(rs.screen_time), 0) as total_screen_time,
				MAX(rs.recorded_at) as last_recorded_at
			FROM group_users gu
			JOIN users u ON u.id = gu.user_id
			JOIN ranked_stats rs ON rs.user_id = gu.user_id AND rs.rn = 1 
			GROUP BY gu.group_id, gu.user_id, u.email, u.username
		)
		SELECT 
			group_id, user_id, total_screen_time,
			AVG(total_screen_time) OVER (PARTITION BY group_id) as group_avg_screen_time
		FROM user_totals
		ORDER BY group_id, total_screen_time DESC;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		startDate,
		endDate,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stats []*WeeklyGroupStats
	for rows.Next() {
		stat := &WeeklyGroupStats{}
		err := rows.Scan(
			&stat.GroupID,
			&stat.UserID,
			&stat.ScreenTime,
			&stat.GroupAverageScreenTime,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil

}
