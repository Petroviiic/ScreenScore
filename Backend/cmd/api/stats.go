package main

import (
	"fmt"
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
	var stats *UserStats

	if err := readJson(w, r, &stats); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "error parsing request body")
		return
	}
	fmt.Println(stats, userId)
	currentRecordTime, err := time.Parse(time.RFC3339, stats.RecordedAt)

	if err != nil {
		fmt.Println("error parsing time")
		return
	}
	fmt.Println("recorded at", currentRecordTime)

	ctx := r.Context()
	if err := app.storage.StatsStorage.AddNewRecord(ctx, int64(userId), stats.ScreenTime, currentRecordTime.UTC()); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// record, err := app.storage.StatsStorage.GetUsersLast(ctx, int64(userId))
	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		//ne dozvoli upis, ispisi poruku
	// 		return
	// 	}
	// }

	// fmt.Println(record.RecordedAt)
}
