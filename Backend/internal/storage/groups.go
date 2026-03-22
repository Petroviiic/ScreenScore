package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
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
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
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

func (g *GroupStorage) JoinGroup(ctx context.Context, userId int64, inviteCode string) (string, error) {
	query := `
		INSERT INTO group_members (group_id, user_id) SELECT id, $1 FROM groups WHERE invite_code = $2 RETURNING group_id;
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var groupId string
	err := g.db.QueryRowContext(
		ctx,
		query,
		userId,
		inviteCode,
	).Scan(
		&groupId,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return "", ERROR_DUPLICATE_KEY_VALUE
			}
		}
		return "", err
	}
	return groupId, err
}

func (g *GroupStorage) LeaveGroup(ctx context.Context, userId int64, groupId string) error {
	query := `
		DELETE FROM group_members WHERE group_id = $1 AND user_id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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

func (g *GroupStorage) KickUser(ctx context.Context, userId int64, groupId string) error {
	return g.LeaveGroup(ctx, userId, groupId)
}

func (g *GroupStorage) GetGroupMembersExclusive(ctx context.Context, groupId string, excludeId int64) ([]int, error) {
	query := `
				SELECT user_id FROM group_members
				WHERE group_id = $1 AND user_id != $2;
				`
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := g.db.QueryContext(
		ctx,
		query,
		groupId,
		excludeId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(
			&id,
		)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil

}
