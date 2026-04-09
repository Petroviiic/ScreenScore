package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type PointsLogicsStorage struct {
	db *sql.DB
}
type WeeklyGroupStats struct {
	GroupID                string
	GroupName              string
	UserID                 int64
	ScreenTime             int
	GroupAverageScreenTime float64
	GroupGoal              float64
	MemberCount            int
	UserRank               int
	PointsToAdd            int
}

func (p *PointsLogicsStorage) GetWeeklyGroupStats(ctx context.Context, weekNumber, weekYear int, startDate, endDate time.Time) ([]*WeeklyGroupStats, error) {
	query := `
		WITH unprocessed_groups AS(	
			SELECT g.id 
			FROM groups g
			WHERE NOT EXISTS (
				SELECT 1 FROM weekly_reward_logs l 
				WHERE l.group_id = g.id 
				AND l.week_number = $1
				AND l.year_year = $2
			)
			LIMIT 100
		),
		group_users AS (
			SELECT user_id, group_id, g.name, g.group_goal FROM group_members 
			JOIN groups g ON g.id = group_id
			WHERE group_id IN (SELECT id FROM unprocessed_groups)
		),
		ranked_stats AS (
			SELECT user_id, screen_time, recorded_at,
				ROW_NUMBER() OVER (PARTITION BY device_id, user_id ORDER BY recorded_at DESC, screen_time DESC) as rn
			FROM screen_time_logs
			WHERE user_id IN (SELECT user_id FROM group_users)
			AND recorded_at >= $3 AND recorded_at < $4
		),
		user_totals AS (
			SELECT 
				gu.name,
				gu.group_goal,
				gu.group_id,
				gu.user_id, 
				u.email, 
				u.username, 
				COALESCE(SUM(rs.screen_time), 0) as total_screen_time,
				MAX(rs.recorded_at) as last_recorded_at
			FROM group_users gu
			JOIN users u ON u.id = gu.user_id
			LEFT JOIN ranked_stats rs ON rs.user_id = gu.user_id AND rs.rn = 1 
			GROUP BY gu.group_id, gu.user_id, u.email, u.username, gu.name, gu.group_goal
		)
		SELECT 
			group_id, 
			group_goal,
			name,
			user_id, 
			total_screen_time, 
			ROW_NUMBER() OVER (PARTITION BY group_id ORDER BY total_screen_time ASC) as user_rank,
			COUNT(*) OVER (PARTITION BY group_id) as members_count,
			AVG(total_screen_time) OVER (PARTITION BY group_id) as group_avg_screen_time
		FROM user_totals
		ORDER BY group_id, total_screen_time ASC;
	`

	rows, err := p.db.QueryContext(
		ctx,
		query,
		weekNumber,
		weekYear,
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
			&stat.GroupGoal,
			&stat.GroupName,
			&stat.UserID,
			&stat.ScreenTime,
			&stat.UserRank,
			&stat.MemberCount,
			&stat.GroupAverageScreenTime,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, nil

}

func (p *PointsLogicsStorage) DistributePoints(ctx context.Context, week, year int, groupRecords map[string][]*WeeklyGroupStats) {
	for groupID, records := range groupRecords {
		err := NewTx(ctx, p.db, func(tx *sql.Tx) error {
			// i := 20
			// for _, record := range records {
			// 	if record.PointsToAdd <= 0 {
			// 		record.PointsToAdd = i

			// 		i += 30
			// 	}
			// }

			if err := batchUpdatePoints(ctx, tx, records); err != nil {
				return err
			}

			if err := batchInsertOnloadNotifications(ctx, tx, records, MESSAGE_WEEKLY_REWARD, NOTIFICATION_TYPE_WEEKLY_REWARD); err != nil {
				return err
			}

			if err := markGroupProcessed(ctx, tx, groupID, week, year); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("Error processing group %s: %v", groupID, err)
		}
	}
}

func batchUpdatePoints(ctx context.Context, tx *sql.Tx, stats []*WeeklyGroupStats) error {
	query := `
		UPDATE users AS u
		SET points = u.points + v.points_to_give
		FROM (VALUES %s) AS v(user_id, points_to_give)
		WHERE u.id = v.user_id;
	`

	valueStrings := make([]string, 0, len(stats))
	valueArgs := make([]interface{}, 0, len(stats)*2)
	for i, s := range stats {
		n := i * 2
		valueStrings = append(valueStrings, fmt.Sprintf("($%d::bigint, $%d::int)", n+1, n+2))
		valueArgs = append(valueArgs, s.UserID, s.PointsToAdd)
	}
	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

	_, err := tx.ExecContext(
		ctx,
		query,
		valueArgs...,
	)
	if err != nil {
		return fmt.Errorf("bulk update failed: %w", err)
	}
	return nil
}
func batchInsertOnloadNotifications(ctx context.Context, tx *sql.Tx, stats []*WeeklyGroupStats, message, notificationType string) error {
	query := `
		INSERT INTO user_notifications 
		(user_id, message, points_earned, is_read, notification_type) 
		VALUES %s;
	`

	valueStrings := make([]string, 0, len(stats))
	valueArgs := make([]interface{}, 0, len(stats)*5)

	i := 0
	for _, s := range stats {
		if s.PointsToAdd <= 0 {
			continue
		}
		n := i * 5
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d::bigint, $%d, $%d::int, $%d::bool, $%d)",
				n+1, n+2, n+3, n+4, n+5))
		valueArgs = append(valueArgs, s.UserID, fmt.Sprintf(message, s.PointsToAdd, s.GroupName), s.PointsToAdd, false, notificationType)

		i++
	}
	if i == 0 {
		return nil
	}
	query = fmt.Sprintf(query, strings.Join(valueStrings, ","))

	_, err := tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("bulk insert failed: %w", err)
	}
	return nil
}

func markGroupProcessed(ctx context.Context, tx *sql.Tx, groupID string, weekNumber, year int) error {
	query := `
		INSERT INTO weekly_reward_logs 
		(group_id, week_number, year_year) 
		VALUES 
		($1, $2, $3);
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		groupID,
		weekNumber,
		year,
	)
	if err != nil {
		return fmt.Errorf("marking group as processed failed: %w", err)
	}
	return nil
}
