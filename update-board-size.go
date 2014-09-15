package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func UpdateBoardSize(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]

		payload struct {
			Width  string `json:"width"`
			Height string `json:"height"`
		}
	)

	// read the payload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// update the board and return the updated board to the user

	board, err := updateBoard(db, user, boardIDHex, bson.M{
		"size": bson.M{
			"width":  payload.Width,
			"height": payload.Height,
		},
	})

	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	utils.WriteJSON(w, http.StatusOK, &board)
	return
}
