package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func UpdateBoard(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]

		payload struct {
			Name        string `json:"name"`
			Background  string `json:"background"`
			Description string `json:"description"`
		}
	)

	board, err := updateBoard(db, user, boardIDHex, bson.M{
		"name":        payload.Name,
		"background":  payload.Background,
		"description": payload.Description,
	})
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	utils.WriteJSON(w, http.StatusOK, &board)
	return
}
