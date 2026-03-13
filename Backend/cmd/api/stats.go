package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type UserStats struct {
	RecordedAt string `json:"recorded_at"`
	ScreenTime int32  `json:"screen_time"`
}

func (app *Application) GetUserStats(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) SyncStats(w http.ResponseWriter, r *http.Request) {
	userId := 1
	var stats UserStats

	if err := readJson(w, r, &stats); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	currentRecordTime, err := time.Parse(time.RFC3339, stats.RecordedAt)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("error parsing time"))
		return
	}
	currentRecordTime = currentRecordTime.UTC()

	if currentRecordTime.After(time.Now().UTC().Add(10 * time.Minute)) {
		log.Println("new record is sent from the future")
		return
	}

	ctx := r.Context()
	lastRecord, err := app.storage.StatsStorage.GetUsersLast(ctx, int64(userId))
	if err != nil {
		if err != sql.ErrNoRows {
			app.badRequestResponse(w, r, err)
			return
		}

		if err := app.storage.StatsStorage.AddNewRecord(ctx, int64(userId), stats.ScreenTime, currentRecordTime); err != nil {
			app.internalServerErrorJson(w, r, err)
			return
		}
		if err := jsonResponse(w, http.StatusAccepted, "database updated"); err != nil {
			app.internalServerErrorJson(w, r, err)
			return
		}
	}

	log.Println("current stats:", stats.ScreenTime, currentRecordTime)
	log.Println("last stats:", lastRecord.ScreenTime, lastRecord.RecordedAt)
	currYear, currMonth, currDay := currentRecordTime.Date()
	lastYear, lastMonth, lastDay := lastRecord.RecordedAt.Date()

	if currentRecordTime.Before(lastRecord.RecordedAt) {
		log.Println("new record time cant be before the last one")
		return
	}

	if currDay == lastDay && currMonth == lastMonth && lastYear == currYear { //same day
		if stats.ScreenTime < lastRecord.ScreenTime {
			return
		}
		if stats.ScreenTime-lastRecord.ScreenTime > int32(currentRecordTime.Sub(lastRecord.RecordedAt).Minutes()) {
			log.Printf("too big screen time difference. the user couldn't have used phone this much in the timespan from the last recorded timestamp")
			return
		}

	}
	todayMidnight := time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, time.UTC)
	if stats.ScreenTime > int32(currentRecordTime.Sub(todayMidnight).Minutes()) {
		log.Printf("current screen time cant be longer than the duration of the current day")
		return
	}

	if err := app.storage.StatsStorage.AddNewRecord(ctx, int64(userId), stats.ScreenTime, currentRecordTime); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusAccepted, "database updated"); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
