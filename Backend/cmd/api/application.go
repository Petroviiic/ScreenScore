package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	_ "github.com/Petroviiic/ScreenScore/docs"
	internalAuth "github.com/Petroviiic/ScreenScore/internal/auth"
	"github.com/Petroviiic/ScreenScore/internal/ratelimiter"
	"github.com/Petroviiic/ScreenScore/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/api/option"
)

type Application struct {
	config           Config
	db               *sql.DB
	storage          *storage.Storage
	authenticator    internalAuth.Authenticator
	rateLimiters     rateLimiters
	firebase         *messaging.Client
	notificationChan chan NotificationTask
}
type rateLimiters struct {
	apiFixedWindow  *ratelimiter.FixedWindowRateLimiter
	authFixedWindow *ratelimiter.FixedWindowRateLimiter
	tokenBucket     *ratelimiter.TokenBucketRatelimiter
}

type Config struct {
	addr            string
	dbConfig        DBConfig
	maxGroupNameLen int
	auth            authConfig
	ratelimiter     rateLimiterConfig
	notifications   notificationsConfig
}

type notificationsConfig struct {
	silentNotificationTimer     time.Duration
	silentNotificationBatchSize int
}
type rateLimiterConfig struct {
	authFixedWindow fixedWindowLimiterConfig
	apiFixedWindow  fixedWindowLimiterConfig
	tokenBucket     tokenBucketLimiterConfig
}
type tokenBucketLimiterConfig struct {
	limit           float64
	tokensPerMinute float64
}
type fixedWindowLimiterConfig struct {
	limit  int
	window time.Duration
}

type authConfig struct {
	secret  string
	expDate time.Duration
	iss     string
}
type DBConfig struct {
	maxIdleConns int
	maxOpenConns int
	maxIdleTime  string
	dbAddr       string
}

type NotificationTask struct {
	UserID int64
	Title  string
	Body   string
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",   // for local dev
			"https://test.vercel.app", // for production
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.GetHealth)

		r.Route("/users", func(r chi.Router) {
			r.Use(app.RatelimiterMiddleware(app.rateLimiters.authFixedWindow, false))

			r.Group(func(r chi.Router) {
				//r.Post("/get-by-id", app.GetById)
				r.Post("/register", app.RegisterUser)
				r.Post("/login", app.Login)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.TokenAuthMiddleware)
				r.Post("/validate_token", app.ValidateJWTToken)
			})
		})

		r.Route("/stats", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)

			r.Group(func(r chi.Router) {
				r.Use(app.RatelimiterMiddleware(app.rateLimiters.tokenBucket, true))
				r.Post("/sync-stats", app.SyncStats)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.RatelimiterMiddleware(app.rateLimiters.apiFixedWindow, true))
				r.Post("/get-group-stats", app.GetGroupStats)
			})
		})

		r.Route("/groups", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)
			r.Use(app.RatelimiterMiddleware(app.rateLimiters.apiFixedWindow, true))

			r.Post("/create/{groupName}", app.CreateGroup)
			r.Post("/join/{inviteCode}", app.JoinGroup)
			r.Post("/leave/{groupId}", app.LeaveGroup)
			r.Post("/kick", app.KickUser)
		})

		r.Route("/notifications", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/sendtestsilent/{token}", app.sendTestNotification)
			})

			r.Group(func(r chi.Router) {
				r.Use(app.TokenAuthMiddleware)

				r.Use(app.RatelimiterMiddleware(app.rateLimiters.apiFixedWindow, true)) //mzd da ovaj custom ne ide na auth ali aj vidjecu
				r.Post("/send_custom", app.SendCustomNotification)
				r.Post("/send_preset", app.SendPresetNotification)
			})
		})
	})

	return r
}
func initFirebase() *messaging.Client {
	opt := option.WithCredentialsFile("screenscore_firebase_key.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Panicf("error initializing firebase app: %v", err)
		return nil
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		log.Panicf("error initializing firebase messaging: %v", err)
		return nil
	}

	return fcmClient
}
func (app *Application) run(router http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("starting server at %s", app.config.addr)

	return srv.ListenAndServe()
}
