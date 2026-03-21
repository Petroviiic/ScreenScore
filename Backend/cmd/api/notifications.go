package main

import (
	"log"
	"net/http"
)

var PresetMessages = map[int]string{
	1: "Touch some grass! 🌿",
	2: "Put the phone down! 📵",
	3: "Go take a nap, screen addict! 😴",
	4: "Is TikTok that interesting? 📱",
	5: "Eyes up, phone down! 🚫",
}

type Notification struct {
	ToUserID int64  `json:"to_user_id" validate:"required"`
	Title    string `json:"title" validate:"required"`
	Body     string `json:"body" validate:"required"`
}

// SendCustomNotification godoc
// @Summary      Sends notification to a user
// @Description  Updates or adds a new screen time record for a specific device.
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        payload  body      Notification  true  "Notification data"
// @Success      200      {string}  string            "Sent"
// @Failure      400      {object}  map[string]string "Bad request payload"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Router       /notifications/send [post]
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

func (app *Application) SendPresetNotification(w http.ResponseWriter, r *http.Request) {
	var payload PresetNotification
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	userID := GetUserFromContext(r)

	if !app.storage.GroupStorage.CheckIfMember(ctx, userID, payload.GroupID) {
		log.Printf("user with id: %d is not a member of group with id: %s", userID, payload.GroupID)
		app.forbiddenResponse(w, r)
		return
	}

	// members, err := app.storage.GroupStorage.GetGroupMembers(ctx, payload.GroupID)
	// if err != nil {
	// 	app.internalServerErrorJson(w, r, err)
	// 	return
	// }
	// for _, val := range members {
	// 	app.notificationChan <- NotificationTask{
	// 		UserID: val.UserID,
	// 		Title:  payload.Title,
	// 		Body:   PresetMessages[payload.MessageId],
	// 	}
	// }

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
