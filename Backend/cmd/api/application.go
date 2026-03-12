package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Application struct {
	config  Config
	db      *sql.DB
	storage *storage.Storage
}

type Config struct {
	addr     string
	dbConfig DBConfig
}

type DBConfig struct {
	maxIdleConns int
	maxOpenConns int
	maxIdleTime  string
	dbAddr       string
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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.GetHealth)

		r.Route("/users", func(r chi.Router) {
			r.Post("/get-by-id", app.GetById)
		})

		r.Route("/stats", func(r chi.Router) {
			r.Post("/sync-stats", app.SyncStats)
			r.Get("/get-stats", app.GetUserStats)
		})
	})

	return r
}

func (app *Application) run(router http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("server started at %s", app.config.addr)
	return srv.ListenAndServe()
}
