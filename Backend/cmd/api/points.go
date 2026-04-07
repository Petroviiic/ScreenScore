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

// GetWeeklyGroupStats godoc
// @Summary      Retrieves weekly group screentime stats including each users total screentime for points calculations
// @Tags         points
// @Produce      json
// @Success      201        {object}  string
// @Failure      403        {object}  map[string]string "Forbidden access"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /points/manual-weekly-reward [get]
func (app *Application) GetWeeklyRewardManually(w http.ResponseWriter, r *http.Request) {
	if app.config.isProdEnv {
		app.forbiddenResponse(w, r)
		return
	}

	groupRecords := app.ProcessWeeklyRewards()

	if err := jsonResponse(w, http.StatusOK, groupRecords); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) ProcessWeeklyRewards() map[string][]*storage.WeeklyGroupStats {
	now := time.Now()
	targetTime := now.AddDate(0, 0, -7)
	year, week := targetTime.ISOWeek()

	startDate, endDate := GetWeekRange(year, week)
	log.Println("week range:", startDate, endDate)

	ctx := context.Background()
	data, err := app.storage.PointsLogicsStorage.GetWeeklyGroupStats(ctx, week, year, startDate, endDate)
	if err != nil {
		log.Println(err)
		return nil
	}

	groupRecords := make(map[string][]*storage.WeeklyGroupStats)
	for _, record := range data {
		groupRecords[record.GroupID] = append(groupRecords[record.GroupID], record)
		if record.GroupAverageScreenTime < float64(app.config.points.MinWeeklyGroupScreentimeThreshold) {
			log.Println("group average screentime is too small. skiping")
			continue
		}
		if record.ScreenTime < app.config.points.MinWeeklyUserScreentimeThreshold {
			log.Printf("user with id %d and group id %s has a screentime of %d. skipping", record.UserID, record.GroupID, record.ScreenTime)
			continue
		}

		diff := record.GroupAverageScreenTime - float64(record.ScreenTime)
		pointsToAdd := int(math.Max(0, math.Round(diff*app.config.points.PointsMultiplier)))

		if record.MemberCount > app.config.points.MinGroupMemberCountThreshold {
			if record.UserRank <= int(math.Ceil(float64(record.MemberCount)*(float64(app.config.points.PercentageOfTopPerformers)/100.0))) {
				pointsToAdd += app.config.points.TopPerformersBonus
			}
		}

		switch record.UserRank {
		case 1:
			pointsToAdd += app.config.points.FirstPlaceBonus
		case 2:
			pointsToAdd += app.config.points.SecondPlaceBonus
		case 3:
			pointsToAdd += app.config.points.ThirdPlaceBonus
		}

		record.PointsToAdd = pointsToAdd
		log.Printf("user with id %d and group id %s gets %d points", record.UserID, record.GroupID, pointsToAdd)

	}

	app.storage.PointsLogicsStorage.DistributePoints(ctx, week, year, groupRecords)

	log.Println("worker processed rewards")
	return groupRecords
}

func (app *Application) PointsWorker() {
	log.Println("Points worker started...")
	ticker := time.NewTicker(app.config.points.PointsTickerTime)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("points worker tick")
				app.ProcessWeeklyRewards()
			default:
				_ = 1
			}
		}
	}()
}

func GetWeekRange(year, week int) (time.Time, time.Time) {
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t, t.AddDate(0, 0, 7)
}
