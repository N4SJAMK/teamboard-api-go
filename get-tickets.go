package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/models"
	"github.com/N4SJAMK/teamboard-api/utils"
)

func GetTickets(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]
	)

	// find the board and make sure the user making the request is a member

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if !board.IsMember(user) {
		utils.Error(w, "User must be a member", http.StatusForbidden)
		return
	}

	// find tickets that are on this board and return them

	var (
		tickets      = []models.Ticket{}
		ticketsQuery = db.C("tickets").Find(bson.M{
			"board_id": board.ID,
		})
	)

	if err := ticketsQuery.All(&tickets); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, tickets)
	return
}
