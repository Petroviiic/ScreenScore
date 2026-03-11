package storage

import (
	"database/sql"
	"fmt"
)

type UserStorage struct {
	db *sql.DB
}

func (u *UserStorage) GetById() {
	fmt.Println("id je taj")
}
