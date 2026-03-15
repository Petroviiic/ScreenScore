package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
)

func isValidString(code string, maxlen int) bool {
	var alphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return len(code) >= 3 && len(code) <= maxlen && alphaNumeric.MatchString(code)
}

type GroupData struct {
	Name       string `json:"name"`
	InviteCode string `json:"invite_code"`
}

func (app *Application) CreateGroup(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "groupName")

	if !isValidString(groupName, app.config.maxGroupNameLen) {
		app.badRequestResponse(w, r, fmt.Errorf("group name malformed"))
		return
	}

	inviteCode, err := app.storage.GroupStorage.CreateGroup(r.Context(), groupName)

	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	group := &GroupData{
		Name:       groupName,
		InviteCode: inviteCode,
	}
	if err := jsonResponse(w, http.StatusCreated, group); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) JoinGroup(w http.ResponseWriter, r *http.Request) {
	userId := int64(1)

	inviteCode := chi.URLParam(r, "inviteCode")

	if !isValidString(inviteCode, 10) {
		app.badRequestResponse(w, r, fmt.Errorf("invite code malformed"))
		return
	}
	err := app.storage.GroupStorage.JoinGroup(r.Context(), userId, inviteCode)

	if err != nil {
		if errors.Is(err, errors.New("no rows affected")) {
			app.internalServerErrorJson(w, r, errors.New("no rows affected"))
			return
		}
		app.internalServerErrorJson(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

func (app *Application) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userId := int64(1)

	groupId := chi.URLParam(r, "groupId")

	err := app.storage.GroupStorage.LeaveGroup(r.Context(), userId, groupId)

	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
