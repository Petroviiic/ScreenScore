package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/messaging"
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
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /notifications/sendtestsilent/{token} [post]
func (app *Application) sendTestNotification(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("Silent notification ticker started...")
	ticker := time.NewTicker(app.config.notifications.silentNotificationTimer)
	i := 0
	return
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("saljem")
				i++
				msg := &messaging.Message{
					Token: "cQ-voKAkQbKNt5zBW__s3z:APA91bGoydj8Jkzx7gMV9SQ1ZI8rKR1Vk0WUi8wenxNOu_SvQfEjiV66jFFkSWFyqpdkoiqvUE7VtAWbKD9rwY_0ScxJ8XwFWZ6rOlEXPWHEMXrvg8Byf3w",
					Notification: &messaging.Notification{
						Title: "samo title",
						Body:  "Marko i ja imamo seks2 " + strconv.Itoa(i),
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
				_, err := app.firebase.Send(context.Background(), msg)

				if err != nil {
					log.Printf("Failed to send notification: %v", err)
				}
				break
				ctx := context.Background()
				tokens, err := app.storage.DeviceStorage.RequestDeviceSync(ctx, app.config.notifications.silentNotificationBatchSize)

				if err != nil {
					log.Println("internal server error, couldn't retrieve devices")
					break
				}
				fmt.Printf("nasao sam %d tokena\n", len(tokens))
				for _, token := range tokens {
					fmt.Println("token: ", token)
					// msg := &messaging.Message{
					// 	Token: token,
					// 	Android: &messaging.AndroidConfig{
					// 		Priority: "high",
					// 	},
					// 	Data: map[string]string{
					// 		"type":       "sync",
					// 		"request_id": "123"},
					// 	APNS: &messaging.APNSConfig{
					// 		Payload: &messaging.APNSPayload{
					// 			Aps: &messaging.Aps{
					// 				ContentAvailable: true,
					// 			},
					// 		},
					// 	},
					// }
					// _, err := app.firebase.Send(ctx, msg)

					// if err != nil {
					// 	log.Printf("Failed to send notification to token %s: %v", token, err)
					// }
				}
			default:
				_ = 1
			}
		}
	}()
}
