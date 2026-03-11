package main

import (
	"log"
	"net/http"
)

func (app *Application) GetHealth(w http.ResponseWriter, r *http.Request) {
	if err := jsonResponse(w, http.StatusOK, "ok"); err != nil {
		log.Panic("something went wrong")
	}
}
