package main

import (
	"net/http"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

type UserPayload struct {
	Username              string `json:"username" validate:"required,max=100"`
	Email                 string `json:"email" validate:"required,email,max=255"`
	Password              string `json:"password" validate:"required,min=3,max=72"`
	DeviceID              string `json:"device_id" validate:"required,max=255"`
	PushNotificationToken string `json:"push_notification_token"`
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
	}
	if err := user.Password.Set(data.Password); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	newUserId, err := app.storage.UserStorage.RegisterUser(r.Context(), user)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if err := app.storage.DeviceStorage.Update(r.Context(), newUserId, data.DeviceID, data.PushNotificationToken); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {

	//...

	// if err := app.storage.DeviceStorage.Update(r.Context(), newUserId, data.DeviceID, data.PushNotificationToken); err != nil {
	// 	app.internalServerErrorJson(w, r, err)
	// 	return
	// }
}
