package main

import (
	"net/http"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

type UserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
	DeviceID string `json:"device_id" validate:"required,max=255"`
}

func (app *Application) GetById(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var data UserPayload
	if err := readJson(w, r, &data); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(data); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &storage.User{
		Email:    data.Email,
		Username: data.Username,
		DeviceID: data.DeviceID,
	}
	if err := user.Password.Set(data.Password); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := app.storage.UserStorage.RegisterUser(r.Context(), user); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
