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

// SyncStreak godoc
// @Summary      Sync streak
// @Tags         streaks
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      403        {object}  map[string]string "User with given id is not a member of the group"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /streak/sync	[post]
func (app *Application) SyncStreak(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserFromContext(r)

	data, err := app.storage.UserStreakStorage.GetStreakData(ctx, userID)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	validation, err := app.ValidateStreak(r.Context(), data, userID, false)

	if err != nil || validation == nil {
		app.customErrorJson(w, r, err, http.StatusExpectationFailed)
		return
	}

	if !validation.StreakFrozen {
		log.Println("streak update")
		if err := app.storage.UserStreakStorage.SaveStreak(ctx, userID, data); err != nil {
			app.internalServerErrorJson(w, r, err)
			return
		}
	} else {
		log.Printf("streak frozen, shields needed %d", validation.ShieldsNeeded)
	}

	response := struct {
		StreakValidation *StreakValidation   `json:"streak_validation"`
		StreakData       *storage.StreakData `json:"streak_data"`
	}{
		StreakValidation: validation,
		StreakData:       data,
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
		lastYear, lastWeek := time.Now().UTC().AddDate(0, 0, -7).ISOWeek()

		if lastWeek != data.WeekNumber || lastYear != data.YearNumber {
			weekStart, weekEnd := GetWeekRange(lastYear, lastWeek)
			screentime, err := app.storage.StatsStorage.GetUserAverageScreenTimeForWeek(ctx, weekStart, weekEnd)

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
			err := fmt.Errorf("not enough shields to repair, have %d, need %d", data.ShieldCount, response.ShieldsNeeded)
			data.CurrentStreak = 0
			data.ShieldCount = 0
			data.LastUpdatedAt = time.Now().UTC()
			return response, err
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
