package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Petroviiic/ScreenScore/internal/storage"
)

var (
	NOT_ENOUGH_SHIELDS_TO_REPAIR_ERROR = fmt.Errorf("not enough shields to repair")
)

type SyncRepairResponse struct {
	StreakValidation *StreakValidation   `json:"streak_validation"`
	StreakData       *storage.StreakData `json:"streak_data"`
}

func (app *Application) StreakLogic(w http.ResponseWriter, r *http.Request, isRepair bool) *SyncRepairResponse {
	ctx := r.Context()
	userID := GetUserFromContext(r)

	data, err := app.storage.UserStreakStorage.GetStreakData(ctx, userID)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return nil
	}

	validation, err := app.ValidateStreak(r.Context(), data, userID, isRepair)
	if err != nil || validation == nil {
		if err != NOT_ENOUGH_SHIELDS_TO_REPAIR_ERROR || validation == nil {
			app.customErrorJson(w, r, err, http.StatusExpectationFailed)
			return nil
		}
	}

	if !validation.StreakFrozen || (validation.StreakFrozen && isRepair) {
		log.Println("streak update")
		if err := app.storage.UserStreakStorage.SaveStreak(ctx, userID, data); err != nil {
			app.internalServerErrorJson(w, r, err)
			return nil
		}
		if validation.StreakFrozen && isRepair {
			validation.StreakLost = true
		}
	} else {
		log.Printf("streak frozen, shields needed %d", validation.ShieldsNeeded)
	}

	response := &SyncRepairResponse{
		StreakValidation: validation,
		StreakData:       data,
	}

	return response
}

// SyncStreak godoc
// @Summary      Sync user streak
// @Description  Calculates the current streak status, checks for inactivity or excessive screen time, and increments the streak if conditions are met.
// @Tags         streaks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  SyncRepairResponse
// @Failure      403  {object}  map[string]string "Unauthorized/Forbidden"
// @Failure      417  {object}  map[string]string "Expectation Failed - Validation error"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /streak/sync [post]
func (app *Application) SyncStreak(w http.ResponseWriter, r *http.Request) {
	response := app.StreakLogic(w, r, false)

	if response == nil {
		return
	}
	if err := jsonResponse(w, http.StatusOK, response); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// RepairStreak godoc
// @Summary      Repair a frozen streak
// @Description  Consumes shields to unfreeze a streak that was blocked due to inactivity or high screen time.
// @Tags         streaks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  SyncRepairResponse
// @Failure      403  {object}  map[string]string "Unauthorized"
// @Failure      417  {object}  map[string]string "Not enough shields or validation failed"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /streak/repair [post]
func (app *Application) RepairStreak(w http.ResponseWriter, r *http.Request) {
	response := app.StreakLogic(w, r, true)

	if response == nil {
		return
	}

	if err := jsonResponse(w, http.StatusOK, response); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

type StreakValidation struct {
	ShieldsNeeded int  `json:"shields_needed"`
	StreakFrozen  bool `json:"streak_frozen"`
	IsNewRecord   bool `json:"is_new_record"`
	StreakLost    bool `json:"streak_lost"`
}

func (app *Application) ValidateStreak(ctx context.Context, data *storage.StreakData, userID int64, repair bool) (*StreakValidation, error) {
	response := &StreakValidation{
		ShieldsNeeded: 0,
		StreakFrozen:  false,
		IsNewRecord:   false,
	}
	if isToday(data.LastUpdatedAt) {
		return response, nil
	}

	today := app.clock.Now().Truncate(24 * time.Hour)
	lastUpdate := data.LastUpdatedAt.UTC().Truncate(24 * time.Hour)
	dayDifference := int(today.Sub(lastUpdate).Hours() / 24)

	yesterdayScreenTime, err := app.storage.StatsStorage.GetUserScreenTimeForDay(ctx, today.AddDate(0, 0, -1), userID)
	if err != nil {
		return nil, err
	}

	if dayDifference > 1 {
		response.ShieldsNeeded = dayDifference - 1
		response.StreakFrozen = true
	}

	if !response.StreakFrozen {
		lastYear, lastWeek := app.clock.Now().AddDate(0, 0, -7).ISOWeek()
		ws, we := GetWeekRange(lastYear, lastWeek)
		fmt.Println(ws, we, yesterdayScreenTime, lastYear, lastWeek)
		if lastWeek != data.WeekNumber || lastYear != data.YearNumber {
			weekStart, weekEnd := GetWeekRange(lastYear, lastWeek)
			screentime, err := app.storage.StatsStorage.GetUserAverageScreenTimeForWeek(ctx, weekStart, weekEnd, userID)

			if err != nil {
				return nil, err
			}
			data.WeekNumber = lastWeek
			data.YearNumber = lastYear
			data.LastWeekAverage = math.Max(app.config.userStreak.minScreenTimeThreshold, screentime)
		}
		if float64(yesterdayScreenTime) > data.LastWeekAverage {
			response.ShieldsNeeded = 1
			response.StreakFrozen = true
		}
		if !response.StreakFrozen {
			data.CurrentStreak++

			if data.CurrentStreak%app.config.userStreak.shieldCountIncreaseRate == 0 {
				if data.ShieldCount < app.config.userStreak.maxShieldCount {
					data.ShieldCount++
				}
			}

			data.LastUpdatedAt = time.Now().UTC()

			if data.CurrentStreak > data.AllTimeHigh {
				data.AllTimeHigh = data.CurrentStreak
				response.IsNewRecord = true
			}
		}
	}

	if repair && response.StreakFrozen {
		if data.ShieldCount >= response.ShieldsNeeded {
			data.ShieldCount -= response.ShieldsNeeded
			response.StreakFrozen = false
			data.LastUpdatedAt = time.Now().UTC()
		} else {
			//fmt.Errorf("not enough shields to repair, have %d, need %d", data.ShieldCount, response.ShieldsNeeded)
			data.CurrentStreak = 0
			data.ShieldCount = 0
			data.LastUpdatedAt = time.Now().UTC()
			return response, NOT_ENOUGH_SHIELDS_TO_REPAIR_ERROR
		}
	}
	return response, nil
}

func isToday(t time.Time) bool {
	now := time.Now().UTC()
	y1, m1, d1 := t.Date()
	y2, m2, d2 := now.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
