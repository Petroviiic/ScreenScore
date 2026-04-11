package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ERROR_NO_ROWS_AFFECTED        = errors.New("no rows affected")
	ERROR_DUPLICATE_KEY_VALUE     = errors.New("record already exists")
	ERROR_NOT_ENOUGH_POINTS_FUNDS = errors.New("not enough funds")
	ERROR_ALREADY_OWN_MESSAGE     = errors.New("you already own the message")
)

type SQLCommon interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Storage struct {
	UserStorage interface {
		GetById(context.Context, int64) (*User, error)
		GetByUsername(ctx context.Context, username string) (*User, error)
		RegisterUser(ctx context.Context, user *User) (int64, error)
		PurchaseMessage(ctx context.Context, messageId int64, userId int64) error
		SpendPoints(ctx context.Context, userId int64, points int) error
		GetPoints(ctx context.Context, userId int64) (int, error)
	}
	StatsStorage interface {
		GetUsersLast(context.Context, int64, string) (*UsageRecord, error)
		AddNewRecord(context.Context, int64, int32, string, time.Time) error
		GetGroupStats(context.Context, string, time.Time) ([]*GroupStats, error)
		GetUserAverageScreenTimeForWeek(context.Context, time.Time, time.Time) (float64, error)
		GetUserScreenTimeForDay(context.Context, time.Time) (int, error)
	}
	GroupStorage interface {
		CheckIfMember(context.Context, int64, string) bool
		CreateGroup(context.Context, string) (string, error)
		GetGroupMembersExclusive(context.Context, string, int64) ([]int, error)
		JoinGroup(ctx context.Context, userId int64, inviteCode string) (string, error)
		LeaveGroup(ctx context.Context, userId int64, groupId string) error
		KickUser(context.Context, int64, string) error
		GetUserGroups(ctx context.Context, userID int64) ([]*Group, error)
		SetGroupGoal(ctx context.Context, goal float64, groupID string) error
	}
	DeviceStorage interface {
		Update(ctx context.Context, userId int64, deviceId string, pushToken string) error
		GetFCMTokens(context.Context, int64) ([]string, error)
		RequestDeviceSync(context.Context, int) ([]string, error)
		UpdateLastSeen(ctx context.Context, userId int64, deviceId string) error
		DeleteFCMToken(string) error
	}
	MessageStorage interface {
		InsertNewPresetMessage(ctx context.Context, tx *sql.Tx, text string, price int, rarity string, isActive bool) error
		GetPresetMessage(context.Context, int64, int64) (string, error)
		GetOwnedPresetMessage(ctx context.Context, userID int64) ([]*PresetMessage, error)
		GetAvaiableInShop(ctx context.Context, userID int64) ([]*PresetMessage, error)
	}
	NotificationStorage interface {
		MarkNotificationAsRead(context.Context, int64, int64) error
		GetUnreadNotifications(ctx context.Context, userID int64) ([]*OnLoadNotification, error)
		AddNewOnLoadNotification(ctx context.Context, userID int64, message string, pointsEarned float64, notificationType string) error
	}
	PointsLogicsStorage interface {
		GetWeeklyGroupStats(ctx context.Context, weekNumber, weekYear int, startDate, endDate time.Time) ([]*WeeklyGroupStats, error)
		DistributePoints(ctx context.Context, week, year int, groupRecords map[string][]*WeeklyGroupStats)
	}
	UserStreakStorage interface {
		GetStreakData(ctx context.Context, userID int64) (*StreakData, error)
		SaveStreak(ctx context.Context, userID int64, data *StreakData) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		UserStorage:         &UserStorage{db},
		StatsStorage:        &StatsStorage{db},
		GroupStorage:        &GroupStorage{db},
		DeviceStorage:       &DeviceStorage{db},
		MessageStorage:      &PresetMessageStorage{db},
		PointsLogicsStorage: &PointsLogicsStorage{db},
		NotificationStorage: &NotificationStorage{db},
		UserStreakStorage:   &UserStreakStorage{db},
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func NewTx(ctx context.Context, db *sql.DB, function func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := function(tx); err != nil {
		return err
	}

	return tx.Commit()
}
