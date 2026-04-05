package storage

import (
	"context"
	"database/sql"
	"time"
)

var (
	NOTIFICATION_TYPE_WEEKLY_REWARD = "weeklyReward"
)
var (
	MESSAGE_WEEKLY_REWARD = "You earned %d points in %s! Check your rank."
)

type NotificationStorage struct {
	db *sql.DB
}
type OnLoadNotification struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	Message          string    `json:"message"`
	PointsEarned     float64   `json:"points_earned"`
	IsRead           bool      `json:"is_read"`
	NotificationType string    `json:"notification_type"`
	CreatedAt        time.Time `json:"created_at"`
}

func (s *NotificationStorage) MarkNotificationAsRead(ctx context.Context, msgID int64, userID int64) error {
	query := `
		UPDATE user_notifications SET is_read = TRUE where id = $1 AND user_id = $2;
	`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res, err := s.db.ExecContext(
		ctx,
		query,
		msgID,
		userID,
	)

	if err != nil {
		return err
	}

	rowNum, _ := res.RowsAffected()
	if rowNum == 0 {
		return ERROR_NO_ROWS_AFFECTED
	}
	return nil
}

func (s *NotificationStorage) GetUnreadNotifications(ctx context.Context, userID int64) ([]*OnLoadNotification, error) {
	query := `	
		SELECT 
			id, 
			user_id, 
			message, 
			points_earned, 
			is_read, 
			notification_type, 
			created_at 
		FROM user_notifications
		WHERE user_id = $1 AND is_read = false;
	`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []*OnLoadNotification
	for rows.Next() {
		notification := &OnLoadNotification{}
		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Message,
			&notification.PointsEarned,
			&notification.IsRead,
			&notification.NotificationType,
			&notification.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *NotificationStorage) AddNewOnLoadNotification(ctx context.Context, userID int64, message string, pointsEarned float64, notificationType string) error {
	query := `
		INSERT INTO user_notifications
		(user_id, message, points_earned, notification_type) 
		VALUES 
		($1, $2, $3, $4);
	`
	_, err := s.db.ExecContext(
		ctx,
		query,
		userID,
		message,
		pointsEarned,
		notificationType,
	)

	if err != nil {
		return err
	}
	return nil

}
