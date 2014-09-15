package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func GetBoard(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]
	)

	// get the board specified by 'board_id'

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	// make sure the user is a member on it

	if !board.IsMember(user) {
		utils.Error(w, "", http.StatusForbidden)
		return
	}

	// return it

	utils.WriteJSON(w, http.StatusOK, board)
	return
}
