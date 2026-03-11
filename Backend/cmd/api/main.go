package main

import (
	"log"

	"github.com/Petroviiic/ScreenScore/internal/db"
	"github.com/Petroviiic/ScreenScore/internal/env"
	"github.com/Petroviiic/ScreenScore/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}
	cfg := Config{
		addr: env.GetString("ADDR", ":3000"),
		dbConfig: DBConfig{
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
			dbAddr:       env.GetString("DB_ADDR", "postgresql://user:user123@localhost:5432/screenscore?sslmode=disable"),
		},
	}

	db, err := db.NewDb(cfg.dbConfig.dbAddr, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleTime)
	if err != nil {
		log.Panic("error connecting to db")
		return
	}

	storage := storage.NewStorage(db)

	app := &Application{
		config:  cfg,
		db:      db,
		storage: storage,
	}

	router := app.mount()

	if err := app.run(router); err != nil {
		log.Panic("error starting the server")
	}
}
