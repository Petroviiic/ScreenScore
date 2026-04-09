package main

import (
	"net/http"

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
	app.ValidateStreak(data)

	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) ValidateStreak(data *storage.StreakData) {

}
