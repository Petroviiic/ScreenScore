package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"firebase.google.com/go/v4/messaging"
	"github.com/Petroviiic/ScreenScore/internal/storage"
	"github.com/go-chi/chi/v5"
)

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
	MessageId int64
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

	ctx := r.Context()
	userID := GetUserFromContext(r)
	msg, err := app.storage.MessageStorage.GetPresetMessage(ctx, userID, payload.MessageId)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if msg == "" {
		app.badRequestResponse(w, r, fmt.Errorf("preset message %d not found", payload.MessageId))
		return
	}

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
	fmt.Println(members, userID)
	points, err := app.storage.UserStorage.GetPoints(ctx, userID)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if points < app.config.notifications.presetNotificationSendingCost {
		app.customErrorJson(w, r, storage.ERROR_NOT_ENOUGH_POINTS_FUNDS, http.StatusBadRequest)
		return
	}
	if err := app.storage.UserStorage.SpendPoints(ctx, userID, app.config.notifications.presetNotificationSendingCost); err != nil {
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

// GetUnreadNotifications godoc
// @Summary      Retrieves a list of all notifications that the users hasn't read
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200      {string}  string            "Notifications"
// @Failure      400      {object}  map[string]string "Bad request payload"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Failure      403      {object}  map[string]string "Forbidden access"
// @Router       /notifications/get_unread [get]
func (app *Application) GetUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	userID := GetUserFromContext(r)
	data, err := app.storage.NotificationStorage.GetUnreadNotifications(r.Context(), userID)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// MarkNotificationAsRead godoc
// @Summary      Notification with provided id will be marked as read
// @Tags         notifications
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        notificationID  path      string  true  "Notification id (URL parameter)"
// @Success      200      {string}  string            "Read"
// @Failure      400      {object}  map[string]string "Bad request payload"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Failure      403      {object}  map[string]string "Forbidden access"
// @Router       /notifications/mark_as_read/{notificationID} [post]
func (app *Application) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	msgIDStr := chi.URLParam(r, "notificationID")
	msgID, err := strconv.ParseInt(msgIDStr, 10, 64)

	userID := GetUserFromContext(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.storage.NotificationStorage.MarkNotificationAsRead(r.Context(), msgID, userID); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, "notification marked as read"); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// GetOwnedMessages godoc
// @Summary      Retrives a list of all preset messages user owns
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200      {string}  string            "List of messages"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Failure      403      {object}  map[string]string "Forbidden access"
// @Router       /notifications/get_users_owned [get]
func (app *Application) GetOwnedMessages(w http.ResponseWriter, r *http.Request) {
	userID := GetUserFromContext(r)
	msgs, err := app.storage.MessageStorage.GetOwnedPresetMessage(r.Context(), userID)

	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, msgs); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// GetAvailableMessagesInShop godoc
// @Summary      Retrives a list of all preset messages user can buy
// @Tags         notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200      {string}  string            "List of messages"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Failure      403      {object}  map[string]string "Forbidden access"
// @Router       /notifications/get_available_shop [get]
func (app *Application) GetAvailableMessagesInShop(w http.ResponseWriter, r *http.Request) {
	userID := GetUserFromContext(r)
	msgs, err := app.storage.MessageStorage.GetAvaiableInShop(r.Context(), userID)

	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, msgs); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// PurchaseMessage godoc
// @Summary      Purchase message from the shop
// @Tags         notifications
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        messageID  path      string  true  "Group id (URL parameter)"
// @Success      200      {string}  string            "Purchased"
// @Failure      400      {object}  map[string]string "Bad request payload"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Failure      403      {object}  map[string]string "Forbidden access"
// @Router       /notifications/purchase/{messageID} [post]
func (app *Application) PurchaseMessage(w http.ResponseWriter, r *http.Request) {
	userID := GetUserFromContext(r)
	msgIDstr := chi.URLParam(r, "messageID")

	msgID, err := strconv.ParseInt(msgIDstr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.storage.UserStorage.PurchaseMessage(r.Context(), msgID, userID); err != nil {
		if errors.Is(err, storage.ERROR_ALREADY_OWN_MESSAGE) {
			app.customErrorJson(w, r, err, http.StatusBadRequest)
			return
		} else if errors.Is(err, storage.ERROR_NOT_ENOUGH_POINTS_FUNDS) {
			app.customErrorJson(w, r, err, http.StatusBadRequest)
			return
		}
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, "Successfully purchased!"); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) StartNotificationWorker() {
	log.Println("Notification worker started...")

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
				if err := app.storage.DeviceStorage.DeleteFCMToken(token); err != nil {
					log.Printf("Token %v is not valid, but couldn't be deleted, error: %v \n", token, err)
				}
			}
		}
	}
}
