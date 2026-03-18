package main

import (
	"log"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/auth"
	"github.com/Petroviiic/ScreenScore/internal/db"
	"github.com/Petroviiic/ScreenScore/internal/env"
	"github.com/Petroviiic/ScreenScore/internal/storage"
	"github.com/joho/godotenv"
)

// @title ScreenScore Backend API
// @version 1.0
// @description API powering ScreenScore backend system.
// @host localhost:3000
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer <your-jwt-token>"
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
		maxGroupNameLen: 15,
		auth: authConfig{
			secret:  env.GetString("AUTH_TOKEN_SECRET", "example"),
			expDate: time.Hour * 24 * 3,
			iss:     env.GetString("AUTH_TOKEN_ISSUER", "admin"),
		},
	}

	db, err := db.NewDb(cfg.dbConfig.dbAddr, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleTime)
	if err != nil {
		log.Panic("error connecting to db")
		return
	}

	storage := storage.NewStorage(db)

	authenticator := auth.NewJWTAuthenticator(cfg.auth.secret, cfg.auth.iss, cfg.auth.iss)
	app := &Application{
		config:        cfg,
		db:            db,
		storage:       storage,
		authenticator: authenticator,
	}

	router := app.mount()

	if err := app.run(router); err != nil {
		log.Panic("error starting the server")
	}
}
