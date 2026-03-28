package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage struct {
	db *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  password  `json:"-"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
}

type password struct {
	Plain string
	Hash  []byte
}

func (p *password) Set(plain string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), 14)
	if err != nil {
		return err
	}

	p.Plain = plain
	p.Hash = hash
	return nil
}

func (p *password) ValidatePassword(plain string) bool {
	if err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plain)); err != nil {
		return false
	}
	return true
}
func (u *UserStorage) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `	SELECT id, email, username, password, points, created_at FROM users 
				WHERE username = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &User{}
	err := u.db.QueryRowContext(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password.Hash,
		&user.Points,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserStorage) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `	SELECT id, email, username, password, points, created_at FROM users 
				WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := &User{}
	err := u.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password.Hash,
		&user.Points,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserStorage) RegisterUser(ctx context.Context, user *User) (int64, error) {
	query := `
			INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id;
		`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userId int64
	err := u.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Username,
		user.Password.Hash,
	).Scan(
		&userId,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return -1, ERROR_DUPLICATE_KEY_VALUE
			}
		}
		return -1, err
	}
	return userId, nil
}

func (m *UserStorage) PurchaseMessage(ctx context.Context, messageId int64, userId int64) error {
	return NewTx(ctx, m.db, func(tx *sql.Tx) error {
		msg, err := getMessageInfo(ctx, tx, messageId)
		if err != nil || msg == nil {
			return err
		}

		points, err := getUserPoints(ctx, tx, userId)
		if err != nil {
			return err
		}

		if msg.Price > points {
			return ERROR_NOT_ENOUGH_POINTS_TO_PURCHASE
		}

		if err := removePoints(ctx, tx, userId, points-msg.Price); err != nil {
			return err
		}

		if err := buyPresetMessage(ctx, tx, userId, messageId); err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					return ERROR_ALREADY_OWN_MESSAGE
				}
			}
			return err
		}
		return nil
	})
}

// helper funcs
func getMessageInfo(ctx context.Context, tx *sql.Tx, messageId int64) (*PresetMessage, error) {
	query := `SELECT id, price, rarity, is_active, created_at FROM preset_messages WHERE id = $1;`

	var msg PresetMessage
	err := tx.QueryRowContext(
		ctx,
		query,
		messageId,
	).Scan(
		&msg.ID,
		&msg.Price,
		&msg.Rarity,
		&msg.IsActive,
		&msg.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
func getUserPoints(ctx context.Context, tx *sql.Tx, userId int64) (int, error) {
	query := `
			SELECT points FROM users WHERE id = $1
			`

	points := -1
	err := tx.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&points,
	)
	if err != nil {
		return -1, err
	}
	return points, nil
}

func removePoints(ctx context.Context, tx *sql.Tx, userId int64, points int) error {
	query := `UPDATE users SET points = $1 WHERE id = $2;`

	_, err := tx.ExecContext(
		ctx,
		query,
		points,
		userId,
	)
	if err != nil {
		return err
	}
	return nil
}

func buyPresetMessage(ctx context.Context, tx *sql.Tx, userId int64, msgID int64) error {
	query := `INSERT INTO user_messages (user_id, message_id) VALUES ($1, $2);`

	_, err := tx.ExecContext(
		ctx,
		query,
		userId,
		msgID,
	)

	if err != nil {
		return err
	}
	return nil
}
