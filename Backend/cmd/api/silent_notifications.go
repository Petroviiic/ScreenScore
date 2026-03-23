package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (app *Application) StartSilentNotificationWorker() {
	fmt.Println("Silent notification ticker started...")
	ticker := time.NewTicker(app.config.notifications.silentNotificationTimer)
	i := 0
	go func() {
		for {
			select {
			case <-ticker.C:
				i++
				// msg := &messaging.Message{
				// 	Token: "c8tfYz4ES7CAGjU2cIlD6h:APA91bExZ1gr90RHalBvQN6O8YFnYDSBSjDO98olq26BjNp5jlyybbMYBCSTQOqlh_aoiAkCXwZz76hnme7jeijo2y4BJVJTYT8OKyl2jfgBBUrAHzDm8ZM",
				// 	Notification: &messaging.Notification{
				// 		Title: "samo title",
				// 		Body:  "Marko i ja imamo seks2 " + strconv.Itoa(i),
				// 	},
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
				// app.firebase.Send(context.Background(), msg)

				// break
				ctx := context.Background()
				tokens, err := app.storage.DeviceStorage.RequestDeviceSync(ctx, app.config.notifications.silentNotificationBatchSize)

				if err != nil {
					log.Println("internal server error, couldn't retrieve devices")
					break
				}
				fmt.Printf("nasao sam %d tokena", len(tokens))
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
