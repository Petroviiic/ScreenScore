package main

import (
	"log"
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

	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) CalculateWeeklyPoints(w http.ResponseWriter, r *http.Request) {

}
