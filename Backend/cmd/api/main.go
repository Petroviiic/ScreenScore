package main

import (
	"log"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/auth"
	"github.com/Petroviiic/ScreenScore/internal/db"
	"github.com/Petroviiic/ScreenScore/internal/env"
	"github.com/Petroviiic/ScreenScore/internal/ratelimiter"
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
		addr:      env.GetString("ADDR", ":3000"),
		isProdEnv: env.GetBool("IS_PROD_ENV", false),
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
		ratelimiter: rateLimiterConfig{
			authFixedWindow: fixedWindowLimiterConfig{
				limit:  15,
				window: 3 * time.Minute,
			},
			apiFixedWindow: fixedWindowLimiterConfig{
				limit:  10,
				window: 1 * time.Minute,
			},
			tokenBucket: tokenBucketLimiterConfig{
				limit:           15,
				tokensPerMinute: 5,
			},
		},
		notifications: notificationsConfig{
			silentNotificationTimer:     30 * time.Minute,
			silentNotificationBatchSize: 100,
			// silentNotificationTimer:     3 * time.Second,
			// silentNotificationBatchSize: 2,
			presetNotificationSendingCost: 5,
		},
		points: pointsConfig{
			MinWeeklyUserScreentimeThreshold:  env.GetInt("MIN_WEEKLY_USER_SCREENTIME_THRESHOLD", 60),
			MinWeeklyGroupScreentimeThreshold: env.GetInt("MIN_WEEKLY_GROUP_SCREENTIME_THRESHOLD", 60),
			MinGroupMemberCountThreshold:      env.GetInt("MIN_GROUP_MEMBER_COUNT_THRESHOLD", 60),
			PercentageOfTopPerformers:         env.GetInt("PERCENTAGE_OF_TOP_PERFORMERS", 30),
			TopPerformersBonus:                env.GetInt("TOP_PERFORMERS_BONUS", 50),
			FirstPlaceBonus:                   env.GetInt("FIRST_PLACE_BONUS", 150),
			SecondPlaceBonus:                  env.GetInt("SECOND_PLACE_BONUS", 100),
			ThirdPlaceBonus:                   env.GetInt("THIRD_PLACE_BONUS", 50),
			GroupGoalBonus:                    env.GetInt("GROUP_GOAL_BONUS", 50),
			PointsTickerTime:                  time.Hour * 3,
			// PointsTickerTime:                  time.Second * 5,
		},
	}

	db, err := db.NewDb(cfg.dbConfig.dbAddr, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleTime)
	if err != nil {
		log.Panic("error connecting to db")
		return
	}
	storage := storage.NewStorage(db)

	authenticator := auth.NewJWTAuthenticator(cfg.auth.secret, cfg.auth.iss, cfg.auth.iss)

	authFixedWindowLimiter := ratelimiter.NewFixedWindowLimiter(cfg.ratelimiter.authFixedWindow.limit, cfg.ratelimiter.authFixedWindow.window)
	authFixedWindowLimiter.Cleanup()

	apiFixedWindowLimiter := ratelimiter.NewFixedWindowLimiter(cfg.ratelimiter.apiFixedWindow.limit, cfg.ratelimiter.apiFixedWindow.window)
	apiFixedWindowLimiter.Cleanup()

	apiTokenBuckerLimiter := ratelimiter.NewTokenBuckerRatelimiter(cfg.ratelimiter.tokenBucket.limit, cfg.ratelimiter.tokenBucket.tokensPerMinute)
	apiTokenBuckerLimiter.Cleanup()

	firebaseApp := initFirebase()
	if firebaseApp == nil {
		log.Panic("firebase client is empty")
	}

	app := &Application{
		config:        cfg,
		db:            db,
		storage:       storage,
		authenticator: authenticator,
		rateLimiters: rateLimiters{
			authFixedWindow: authFixedWindowLimiter,
			apiFixedWindow:  apiFixedWindowLimiter,
			tokenBucket:     apiTokenBuckerLimiter,
		},
		firebase:         firebaseApp,
		notificationChan: make(chan NotificationTask),
	}

	router := app.mount()

	go app.StartNotificationWorker()
	app.StartSilentNotificationWorker()
	app.PointsWorker()

	if err := app.run(router); err != nil {
		log.Panic("error starting the server")
	}
}
