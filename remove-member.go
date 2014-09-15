package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func RemoveMember(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex  = c.URLParams["board_id"]
		memberIDHex = c.URLParams["member_id"]
	)

	// get the board for this request, the user must be an 'admin'

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if board.GetRole(user) != "admin" {
		utils.Error(w, "", http.StatusForbidden)
		return
	}

	// make sure the 'user_id' is a valid ObjectID, and that
	// the member actually exists

	if !bson.IsObjectIdHex(memberIDHex) {
		utils.Error(w, "'user_id' must be a valid 'ObjectID'",
			http.StatusBadRequest)
		return
	}
	memberID := bson.ObjectIdHex(memberIDHex)

	member := board.GetMember(memberID)
	if member == nil {
		utils.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	// remove the membership and return the removed member

	if err := db.C("boards").UpdateId(board.ID, bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"user_id": member.UserID,
			},
		},
	}); err != nil {
		if err == mgo.ErrNotFound {
			utils.Error(w, "Board not found", http.StatusNotFound)
			return
		}
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, member)
	return
}
