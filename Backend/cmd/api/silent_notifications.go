package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/go-chi/chi/v5"
)

// SendTestSilentNotification godoc
// @Summary      Sends a test silent notification
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        token  path      string  true  "Token (URL parameter)"
// @Success      201        {object}  string
// @Failure      400        {object}  map[string]string "Token malformed"
// @Failure      403        {object}  map[string]string "Forbidden access"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /notifications/sendtestsilent/{token} [post]
func (app *Application) sendTestNotification(w http.ResponseWriter, r *http.Request) {
	if app.config.isProdEnv {
		app.forbiddenResponse(w, r)
		return
	}
	token := chi.URLParam(r, "token")

	if token == "" {
		app.badRequestResponse(w, r, fmt.Errorf("token is empty"))
		return
	}

	token = strings.ReplaceAll(token, "%3A", ":")
	msg := &messaging.Message{
		Token: token,
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
	_, err := app.firebase.Send(r.Context(), msg)

	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, "notification sent"); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
func (app *Application) StartSilentNotificationWorker() {
	log.Println("Silent notification ticker started...")
	ticker := time.NewTicker(app.config.notifications.silentNotificationTimer)
	//i := 0
	//return
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("requesting data sync")
				ctx := context.Background()
				tokens, err := app.storage.DeviceStorage.RequestDeviceSync(ctx, app.config.notifications.silentNotificationBatchSize)

				if err != nil {
					log.Println("internal server error, couldn't retrieve devices")
					break
				}
				if len(tokens) == 0 {
					log.Println("no tokens found, waiting for next sync")
					continue
				}
				numBatches := (len(tokens) + app.config.notifications.silentNotificationBatchSize - 1) / app.config.notifications.silentNotificationBatchSize

				go func() {
					for i := 0; i < numBatches; i++ {
						var messages []*messaging.Message
						for j := 0; j < app.config.notifications.silentNotificationBatchSize; j++ {
							index := i*app.config.notifications.silentNotificationBatchSize + j
							if index >= len(tokens) {
								log.Println("index > len(tokens), breaking")
								break
							}
							token := tokens[index]

							messages = append(messages, app.GetSilentMessage(ctx, token))
						}

						batchResponse, err := app.firebase.SendEach(ctx, messages)
						if err != nil {
							log.Printf("Batch %d failed: %v", i, err)
						} else {
							log.Printf("Batch %d: Success %d, Failure %d", i, batchResponse.SuccessCount, batchResponse.FailureCount)
						}

						for idx, resp := range batchResponse.Responses {
							if !resp.Success {
								log.Println("!resp success, ", idx)
								if messaging.IsUnregistered(resp.Error) {
									if err := app.storage.DeviceStorage.DeleteFCMToken(messages[idx].Token); err != nil {
										log.Printf("Token %v is not valid, but couldn't be deleted, error: %v \n", messages[idx].Token, err)
									}
								}
							}
						}
						time.Sleep(app.config.notifications.silentNotificationTimer / time.Duration(numBatches))
					}

				}()
			default:
				_ = 1
			}
		}
	}()
}

func (app *Application) GetSilentMessage(ctx context.Context, token string) *messaging.Message {
	msg := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"type":       "sync",
			"request_id": "123"},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
			},
		},
	}
	return msg
}
