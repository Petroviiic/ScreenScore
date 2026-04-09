package storage

import (
	"context"
	"database/sql"
	"time"
)

type UserStreakStorage struct {
	db *sql.DB
}

type StreakData struct {
	ID              int64     `json:"id"`
	CurrentStreak   int       `json:"current_streak"`
	AllTimeHigh     int       `json:"all_time_high"`
	ShieldCount     int       `json:"shield_count"`
	WeekNumber      int       `json:"week_number"`
	YearNumber      int       `json:"year_number"`
	LastWeekAverage float64   `json:"last_week_average"`
	LastUpdatedAt   time.Time `json:"last_updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

func (s *UserStreakStorage) GetStreakData(ctx context.Context, userID int64) (*StreakData, error) {
	query := `
		SELECT 
			id, 
			current_streak, 
			all_time_high, 
			shield_count, 
			week_number, 
			year_number, 
			last_week_average, 
			last_updated_at, 
			created_at 
		FROM user_streaks
		WHERE user_id = $1
		;`

	now := time.Now().UTC()
	lastWeek, lastYear := now.AddDate(0, 0, -7).ISOWeek()
	data := &StreakData{
		CurrentStreak:   0,
		AllTimeHigh:     0,
		ShieldCount:     0,
		WeekNumber:      lastWeek,
		YearNumber:      lastYear,
		LastWeekAverage: 200,
		LastUpdatedAt:   time.Now().AddDate(0, 0, -1).UTC(),
	}

	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&data.ID,
		&data.CurrentStreak,
		&data.AllTimeHigh,
		&data.ShieldCount,
		&data.WeekNumber,
		&data.YearNumber,
		&data.LastWeekAverage,
		&data.LastUpdatedAt,
		&data.CreatedAt,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return data, nil
}

func (s *UserStreakStorage) SaveStreak(ctx context.Context, data *StreakData) error {
	// query := `
	// 		INSERT INTO user_streaks (user_id, current_streak, all_time_high, shield_count, week_number, year_number, last_week_average, last_updated_at)
	// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	// 		ON CONFLICT (user_id)
	// 		DO UPDATE SET
	// 			current_streak = EXCLUDED.current_streak,
	// 			all_time_high = GREATEST(user_streaks.all_time_high, EXCLUDED.all_time_high),
	// 			shield_count = EXCLUDED.shield_count,
	// 			last_updated_at = EXCLUDED.last_updated_at;
	//
	//		--returning id --vidi da li ti treba ovaj
	// 	`

	return nil
}
