package main

import "net/http"

type Notification struct {
	UserID int64  `json:"user_id" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

// SendNotification godoc
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
func (app *Application) SendNotification(w http.ResponseWriter, r *http.Request) {
	var payload Notification
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.notificationChan <- NotificationTask{
		UserID: payload.UserID,
		Title:  payload.Title,
		Body:   payload.Body,
	}

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
