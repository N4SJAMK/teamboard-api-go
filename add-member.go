package main

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/models"
	"github.com/N4SJAMK/teamboard-api/utils"
)

func AddMember(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]

		payload struct {
			Role   string        `json:"role"`
			UserID bson.ObjectId `json:"user_id"`
		}
	)

	// read the payload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate that the 'role' in payload matches either 'admin' or 'member'

	if payload.Role != "admin" && payload.Role != "member" {
		utils.Error(w, "'role' must be either 'admin' or 'member'",
			http.StatusBadRequest)
		return
	}

	// get the requested board

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	// make sure the user making the request is 'admin' on the board

	if board.GetRole(user) != "admin" {
		utils.Error(w, "Role 'admin' required", http.StatusForbidden)
		return
	}

	// make sure the user we are about to add is not already a member
	// and that the user actually exists

	if board.IsMember(payload.UserID) {
		utils.Error(w, "User is already a member", http.StatusConflict)
		return
	}

	// TODO what is the most efficient way of checking that the
	//      document we are referencing with 'user_id' actually exists

	// add the user as a member

	newMember := models.Member{
		Role:        payload.Role,
		UserID:      payload.UserID,
		MemberSince: time.Now(),
	}

	if err := db.C("boards").UpdateId(board.ID, bson.M{
		"$push": bson.M{
			"members": newMember,
		},
	}); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the member we just added

	utils.WriteJSON(w, http.StatusCreated, &newMember)
	return
}
