package main

import (
	"log"
	"math"
	"net/http"
	"time"
)

// GetWeeklyGroupStats godoc
// @Summary      Retrieves weekly group screentime stats including each users total screentime for points calculations
// @Tags         points
// @Produce      json
// @Success      201        {object}  string
// @Failure      403        {object}  map[string]string "Forbidden access"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /points/get_group_stats [get]
func (app *Application) GetWeeklyGroupStats(w http.ResponseWriter, r *http.Request) {
	if app.config.isProdEnv {
		app.forbiddenResponse(w, r)
		return
	}

	eYear, eMonth, eDay := time.Now().UTC().Date()
	endDate := time.Date(eYear, eMonth, eDay, 0, 0, 0, 0, time.UTC)
	// startDate := endDate.AddDate(0, 0, -7)
	startDate := endDate.AddDate(0, 0, -20) //TODO delete this

	log.Println(startDate, endDate)

	data, err := app.storage.PointsLogicsStorage.GetWeeklyGroupStats(r.Context(), startDate, endDate)
	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	for _, record := range data {
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
			if record.UserRank <= int(math.Ceil(float64(record.MemberCount*app.config.points.PercentageOfTopPerformers/100))) {
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

	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) CalculateWeeklyPoints(w http.ResponseWriter, r *http.Request) {

}
