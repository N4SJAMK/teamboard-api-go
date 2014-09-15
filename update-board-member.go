package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func UpdateBoardMember(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex  = c.URLParams["board_id"]
		memberIDHex = c.URLParams["member_id"]

		payload struct {
			Role string `json:"height"`
		}
	)

	// read the payload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// make sure the 'member_id' attribute is a valid 'ObjectID'

	if !bson.IsObjectIdHex(memberIDHex) {
		utils.Error(w, "'user_id' must be a valid 'ObjectID'",
			http.StatusBadRequest)
		return
	}
	memberID := bson.ObjectIdHex(memberIDHex)

	// get the board and make sure the user has admin access

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if board.GetRole(user) != "admin" {
		utils.Error(w, "", http.StatusForbidden)
		return
	}

	// make sure the 'member_id' is an actual member

	member := board.GetMember(memberID)
	if member == nil {
		utils.Error(w, "Member not found on Board", http.StatusNotFound)
		return
	}

	// apply updates and return the newly updated board
	// NOTE this does not use the 'updateBoard' helper function
	//      because this uses the positional operator to update
	//      the embedded member document

	var (
		boardQuery = bson.M{
			"_id":             board.ID,
			"members.user_id": member.UserID,
		}
		boardUpdate = bson.M{
			"$set": bson.M{
				"members.$.role": payload.Role,
			},
		}
	)

	if err := db.C("boards").Update(boardQuery, boardUpdate); err != nil {
		if err == mgo.ErrNotFound {
			utils.Error(w, "Board not found", http.StatusNotFound)
			return
		}
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	member.Role = payload.Role
	utils.WriteJSON(w, http.StatusOK, &member)
	return
}
