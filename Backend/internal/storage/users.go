package storage

import (
	"context"
	"database/sql"
	"time"
)

type UserStorage struct {
	db *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *UserStorage) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `	SELECT id, email, username, password, created_at FROM users 
				WHERE id = $1`

	user := &User{}
	err := u.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
