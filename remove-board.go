package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func RemoveBoard(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]
	)

	// get the board specified by 'board_id' and make sure the user
	// is an 'admin' on the board

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if board.GetRole(user) != "admin" {
		utils.Error(w, "User must have 'admin' privileges",
			http.StatusForbidden)
		return
	}

	// remove board, and the tickets that are associated with it

	if err := db.C("boards").RemoveId(board.ID); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := db.C("tickets").RemoveAll(bson.M{
		"board_id": board.ID,
	}); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the removed board

	utils.WriteJSON(w, http.StatusOK, board)
	return
}
