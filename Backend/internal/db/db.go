package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDb(addr string, maxIdleConns, maxOpenConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)

	if err != nil {
		log.Println("JEBEM OVDJE")
		return nil, err
	}

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Println("JEBEM OVDJE 2")
		return nil, err
	}

	return db, nil
}
