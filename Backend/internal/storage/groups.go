package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type GroupsStorage struct {
	db *sql.DB
}

func (g *GroupsStorage) CheckIfMember(ctx context.Context, userId int64, groupId string) bool {
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
