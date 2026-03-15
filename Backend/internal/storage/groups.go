package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type GroupStorage struct {
	db *sql.DB
}

func (g *GroupStorage) CheckIfMember(ctx context.Context, userId int64, groupId string) bool {
	query := `	SELECT EXISTS (
					SELECT 1 FROM group_members 
					WHERE group_id = $1 AND user_id = $2
				);
			`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	isMember := false
	err := g.db.QueryRowContext(
		ctx,
		query,
		groupId,
		userId,
	).Scan(
		&isMember,
	)

	if err != nil {
		fmt.Println(err)
		return false
	}
	return isMember
}

func RandomInviteCode(desiredLen int) string {
	b := make([]byte, desiredLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func (g *GroupStorage) CreateGroup(ctx context.Context, groupName string) (string, error) {
	query := `
				INSERT INTO groups (name, invite_code) VALUES ($1, $2);
			`

	inviteCode := RandomInviteCode(10)

	_, err := g.db.ExecContext(
		ctx,
		query,
		groupName,
		inviteCode,
	)

	if err != nil {
		return "", err
	}

	return inviteCode, nil
}

func (g *GroupStorage) JoinGroup(ctx context.Context, userId int64, inviteCode string) error {
	query := `
		INSERT INTO group_members (group_id, user_id) SELECT id, $1 FROM groups WHERE invite_code = $2;
	`
	resp, err := g.db.ExecContext(
		ctx,
		query,
		userId,
		inviteCode,
	)
	if err != nil {
		return err
	}
	num, _ := resp.RowsAffected()
	if num == 0 {
		return errors.New(ERROR_NO_ROWS_AFFECTED)
	}
	return nil
}

func (g *GroupStorage) LeaveGroup(ctx context.Context, userId int64, groupId string) error {
	query := `
		DELETE FROM group_members WHERE group_id = $1 AND user_id = $2;
	`
	_, err := g.db.ExecContext(
		ctx,
		query,
		groupId,
		userId,
	)
	if err != nil {
		return err
	}
	return nil
}
