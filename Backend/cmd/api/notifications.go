package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/messaging"
)

var PresetMessages = map[int]string{
	1: "Touch some grass!",
	2: "Put the phone down!",
	3: "Go take a nap, screen addict!",
	4: "Is TikTok that interesting?",
	5: "Eyes up, phone down!",
}

type Notification struct {
	ToUserID int64  `json:"to_user_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
}

// SendCustomNotification godoc
// @Summary      Sends a custom notification to a user
// @Tags         notifications
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      Notification  true  "Notification data"
// @Success      200      {string}  string            "Sent"
// @Failure      400      {object}  map[string]string "Bad request payload"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Router       /notifications/send_custom [post]
func (app *Application) SendCustomNotification(w http.ResponseWriter, r *http.Request) {
	var payload Notification
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.notificationChan <- NotificationTask{
		UserID: payload.ToUserID,
		Title:  payload.Title,
		Body:   payload.Body,
	}

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

type PresetNotification struct {
	GroupID   string
	MessageId int
}

// SendPresetNotification godoc
// @Summary      Sends preset notification to a group
// @Tags         notifications
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      PresetNotification  true  "PresetNotification data"
// @Success      200      {string}  string            "Sent"
// @Failure      400      {object}  map[string]string "Bad request payload"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Failure      403      {object}  map[string]string "Forbidden access"
// @Router       /notifications/send_preset [post]
func (app *Application) SendPresetNotification(w http.ResponseWriter, r *http.Request) {
	var payload PresetNotification
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	msg, ok := PresetMessages[payload.MessageId]
	if !ok {
		app.badRequestResponse(w, r, fmt.Errorf("preset message %d not found", payload.MessageId))
		return
	}

	ctx := r.Context()
	userID := GetUserFromContext(r)

	if !app.storage.GroupStorage.CheckIfMember(ctx, userID, payload.GroupID) {
		log.Printf("user with id: %d is not a member of group with id: %s", userID, payload.GroupID)
		app.forbiddenResponse(w, r)
		return
	}
	user, err := app.storage.UserStorage.GetById(ctx, userID)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	members, err := app.storage.GroupStorage.GetGroupMembersExclusive(ctx, payload.GroupID, userID)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	for _, val := range members {
		app.notificationChan <- NotificationTask{
			UserID: int64(val),
			Title:  fmt.Sprintf("%s says: ", user.Username),
			Body:   msg,
		}
	}

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) StartNotificationWorker() {
	log.Println("Notification worker started...")

	msg := &messaging.Message{
		Token:        "f18auD8cQ--YUIfy9cfdb6:APA91bF4yXNXLQsSTGdeCSGPpUP97Rr1lHerSXnlAEZTsxTTBA9Q2rmMQBnDNsyapF1_Nb-i2pnXhIyNOkUfHLiys4ZRBcFHPXdZs2PnBN6FwKeb-kKU4SI",
		Notification: &messaging.Notification{
			//Title: "samo title",
			//Body: "Marko i ja imamo seks",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		Data: map[string]string{
			"type":       "sync",
			"request_id": "123"},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
			},
		},
	}
	app.firebase.Send(context.Background(), msg)

	for task := range app.notificationChan {
		ctx := context.Background()

		tokens, err := app.storage.DeviceStorage.GetFCMTokens(ctx, task.UserID)

		if err != nil {
			log.Printf("Could not get tokens for user %d: %v", task.UserID, err)
			continue
		}

		for _, token := range tokens {
			msg := &messaging.Message{
				Token: token,
				Notification: &messaging.Notification{
					Title: task.Title,
					Body:  task.Body,
				},
			}

			_, err := app.firebase.Send(ctx, msg)
			if err != nil {
				log.Printf("Failed to send notification to token %s: %v", token, err)
			}
		}
	}
}
