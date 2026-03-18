package main

import (
	"fmt"
	"log"
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

// CreateGroup godoc
// @Summary      Create new group
// @Description  Creates a new group and returns invite code
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        groupName  path      string  true  "Group name (URL parameter)"
// @Success      201        {object}  GroupData
// @Failure      400        {object}  map[string]string "Group name malformed"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /groups/create/{groupName} [post]
func (app *Application) CreateGroup(w http.ResponseWriter, r *http.Request) {
	//TODO add auth to docs

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

// JoinGroup godoc
// @Summary      Join group
// @Description  Joins a group and returns group id
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        inviteCode  path      string  true  "Invite code (URL parameter)"
// @Success      200		{object}  string
// @Failure      400        {object}  map[string]string "Invite code malformed"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /groups/join/{inviteCode} [post]
func (app *Application) JoinGroup(w http.ResponseWriter, r *http.Request) {
	userId := GetUserFromContext(r)

	inviteCode := chi.URLParam(r, "inviteCode")

	if !isValidString(inviteCode, 10) {
		app.badRequestResponse(w, r, fmt.Errorf("invite code malformed"))
		return
	}
	groupId, err := app.storage.GroupStorage.JoinGroup(r.Context(), userId, inviteCode)

	if err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, groupId); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}

// LeaveGroup godoc
// @Summary      Leave group
// @Description  Leaves a group
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        groupId  path      string  true  "Group id (URL parameter)"
// @Success      200
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /groups/leave	[post]
func (app *Application) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userId := GetUserFromContext(r)

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

type KickUserPayload struct {
	UserToKickID int64  `json:"user_to_kick_id" validate:"required"`
	GroupID      string `json:"group_id" validate:"required"`
}

// KickUser godoc
// @Summary      Kick user
// @Description  Anyone can kick another group member
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      KickUserPayload  true  "Payload with user to kick and group id"
// @Success      200
// @Failure      403        {object}  map[string]string "User with given id is not a member of the group"
// @Failure      500        {object}  map[string]string "Internal server error"
// @Router       /groups/kick	[post]
func (app *Application) KickUser(w http.ResponseWriter, r *http.Request) {
	//TODO add auth to docs

	var payload KickUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if !app.storage.GroupStorage.CheckIfMember(ctx, payload.UserToKickID, payload.GroupID) {
		log.Printf("user with id: %d is not a member of group with id: %s", payload.UserToKickID, payload.GroupID)
		app.forbiddenResponse(w, r)
		return
	}
	if err := app.storage.GroupStorage.KickUser(ctx, payload.UserToKickID, payload.GroupID); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerErrorJson(w, r, err)
		return
	}
}
