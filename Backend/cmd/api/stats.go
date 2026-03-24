package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

var _ = storage.GroupStats{}

type GroupStatsPayload struct {
	GroupID     string `json:"group_id" validate:"required"`
	DesiredDate string `json:"desired_date" validate:"required" example:"2026-03-17"`
}

// GetGroupStats godoc
// @Summary      Retrieves screentime data for all group memebers
// @Tags         stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      GroupStatsPayload  true  "Group stats payload"
// @Success      200		{array}   storage.GroupStats
// @Failure      403        {object}  map[string]string "User with given id is not a member of the group"
// @Failure      400        {object}  map[string]string "Payload malformed"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /stats/get-group-stats	[post]
func (app *Application) GetGroupStats(w http.ResponseWriter, r *http.Request) {
	userId := GetUserFromContext(r)

	var payload GroupStatsPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if !app.storage.GroupStorage.CheckIfMember(ctx, userId, payload.GroupID) {
		log.Printf("user with id: %d is not a member of group with id: %s", userId, payload.GroupID)
		app.forbiddenResponse(w, r)
		return
	}

	desiredDate, _ := time.Parse(time.DateOnly, payload.DesiredDate)

	stats, err := app.storage.StatsStorage.GetGroupStats(ctx, payload.GroupID, desiredDate)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, stats); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

type UserStatsPayload struct {
	DeviceID   string `json:"device_id" validate:"required"`
	RecordedAt string `json:"recorded_at" validate:"required" example:"2026-03-17T12:00:00Z"`
	ScreenTime int32  `json:"screen_time" validate:"required"`
}

// SyncStats godoc
// @Summary      Synchronize user screen time
// @Description  Updates or adds a new screen time record for a specific device.
// @Description  Includes logic to prevent "cheating" (time from future, screen time increasing faster than real time, etc.)
// @Tags         stats
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      UserStatsPayload  true  "User screen time data"
// @Success      201      {string}  string            "database updated"
// @Failure      400      {object}  map[string]string "Validation error (Future time, time travel, invalid increments)"
// @Failure      500      {object}  map[string]string "Internal server error"
// @Router       /stats/sync-stats [post]
func (app *Application) SyncStats(w http.ResponseWriter, r *http.Request) {
	userId := GetUserFromContext(r)
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
	log.Println(currentRecordTime, time.Now().UTC().Add(10*time.Minute))
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
	_, offset := currentRecordTime.Local().Zone()
	zoneOffsetMinutes := int32(offset / 60)

	todayMidnight := time.Date(currYear, currMonth, currDay, 0, 0, 0, 0, time.UTC)
	if (stats.ScreenTime - zoneOffsetMinutes) > int32(currentRecordTime.Sub(todayMidnight).Minutes()) {
		app.badRequestResponse(w, r, fmt.Errorf("current screen time cant be longer than the duration of the current day"))
		return
	}

	if err := app.storage.StatsStorage.AddNewRecord(ctx, int64(userId), stats.ScreenTime, stats.DeviceID, currentRecordTime); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := app.storage.DeviceStorage.UpdateLastSeen(ctx, userId, stats.DeviceID); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusCreated, "database updated"); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
