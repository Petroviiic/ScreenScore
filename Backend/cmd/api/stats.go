package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
	"github.com/go-chi/chi/v5"
)

var _ = storage.GroupStats{}

type UserStatsPayload struct {
	DeviceID   string `json:"device_id" validate:"required"`
	RecordedAt string `json:"recorded_at" validate:"required" example:"2026-03-17T12:00:00Z"`
	ScreenTime int32  `json:"screen_time" validate:"required"`
}

// GetGroupStats godoc
// @Summary      Retrieves screentime data for all group memebers
// @Tags         stats
// @Accept       json
// @Produce      json
// @Param        groupID  path      string  true  "Group id (URL parameter)"
// @Success      200		{array}   storage.GroupStats
// @Failure      403        {object}  map[string]string "User with given id is not a member of the group"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /stats/get-group-stats/{groupID}	[get]
func (app *Application) GetGroupStats(w http.ResponseWriter, r *http.Request) {
	//TODO - userId, add auth
	userId := int64(2)

	groupId := chi.URLParam(r, "groupID")
	ctx := r.Context()
	if !app.storage.GroupStorage.CheckIfMember(ctx, userId, groupId) {
		log.Printf("user with id: %d is not a member of group with id: %s", userId, groupId)
		app.forbiddenResponse(w, r)
		return
	}

	stats, err := app.storage.StatsStorage.GetGroupStats(ctx, groupId)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, stats); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// SyncStats godoc
// @Summary      Synchronize user screen time
// @Description  Updates or adds a new screen time record for a specific device.
// @Description  Includes logic to prevent "cheating" (time from future, screen time increasing faster than real time, etc.)
// @Tags         stats
// @Accept       json
// @Produce      json
// @Param        payload  body      UserStatsPayload  true  "User screen time data"
// @Success      201      {string}  string            "database updated"
// @Failure      400      {object}  map[string]string "Validation error (Future time, time travel, invalid increments)"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Router       /stats/sync-stats [post]
func (app *Application) SyncStats(w http.ResponseWriter, r *http.Request) {
	userId := 1 //TODO, add auth docs
	var stats UserStatsPayload

	if err := readJson(w, r, &stats); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(stats); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	currentRecordTime, err := time.Parse(time.RFC3339, stats.RecordedAt)
	if err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("error parsing time"))
		return
	}
	currentRecordTime = currentRecordTime.UTC()

	if currentRecordTime.After(time.Now().UTC().Add(10 * time.Minute)) {
		app.badRequestResponse(w, r, fmt.Errorf("new record is sent from the future"))
		return
	}

	ctx := r.Context()
	lastRecord, err := app.storage.StatsStorage.GetUsersLast(ctx, int64(userId), stats.DeviceID)
	if err != nil {
		if err != sql.ErrNoRows {
			app.badRequestResponse(w, r, err)
			return
		}

		if err := app.storage.StatsStorage.AddNewRecord(ctx, int64(userId), stats.ScreenTime, stats.DeviceID, currentRecordTime); err != nil {
			app.internalServerErrorJson(w, r, err)
			return
		}
		if err := jsonResponse(w, http.StatusCreated, "database updated"); err != nil {
			app.internalServerErrorJson(w, r, err)
			return
		}
		return
	}

	// log.Println("current stats:", stats.ScreenTime, currentRecordTime)
	// log.Println("last stats:", lastRecord.ScreenTime, lastRecord.RecordedAt)
	currYear, currMonth, currDay := currentRecordTime.Date()
	lastYear, lastMonth, lastDay := lastRecord.RecordedAt.Date()

	if currentRecordTime.Before(lastRecord.RecordedAt) {
		app.badRequestResponse(w, r, fmt.Errorf("new record timestamp cannot be earlier than the last one"))
		return
	}

	if currDay == lastDay && currMonth == lastMonth && lastYear == currYear { //same day
		if stats.ScreenTime < lastRecord.ScreenTime {
			app.badRequestResponse(w, r, fmt.Errorf("new screen time cannot be lower than the previous record"))
			return
		}
		if stats.ScreenTime-lastRecord.ScreenTime > int32(currentRecordTime.Sub(lastRecord.RecordedAt).Minutes()) {
			app.badRequestResponse(w, r, fmt.Errorf("screen time increase exceeds elapsed real time"))
			return
		}
		if stats.ScreenTime == lastRecord.ScreenTime {
			app.badRequestResponse(w, r, fmt.Errorf("screen time is the same as the last one"))
			return
		}

	}
	todayMidnight := time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, time.UTC)
	if stats.ScreenTime > int32(currentRecordTime.Sub(todayMidnight).Minutes()) {
		app.badRequestResponse(w, r, fmt.Errorf("current screen time cant be longer than the duration of the current day"))
		return
	}

	if err := app.storage.StatsStorage.AddNewRecord(ctx, int64(userId), stats.ScreenTime, stats.DeviceID, currentRecordTime); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, "database updated"); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
