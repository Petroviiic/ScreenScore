package main

import (
	"log"

	"github.com/Petroviiic/ScreenScore/internal/env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading env file: %v", err)
	}
	cfg := Config{
		dbConfig: DBConfig{
			username: env.GetString("DB_USER", "user"),
			password: env.GetString("DB_PASSWORD", "user123"),
			dbName:   env.GetString("DB_NAME", "screenscore"),
			dbHost:   env.GetString("DB_HOST", "localhost"),
			dbAddr:   env.GetString("DB_ADDR", "postgresql://user:user123@localhost:5432/screenscore?sslmode=disable"),
		},
	}

	app := &Application{
		config: cfg,
	}

	app.Mount()
}
