package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	config Config
}

type Config struct {
	dbConfig DBConfig
}

type DBConfig struct {
	username string
	password string
	dbName   string
	dbHost   string
	dbAddr   string
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

	})
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.GetHealth)
	})

	return r
}
