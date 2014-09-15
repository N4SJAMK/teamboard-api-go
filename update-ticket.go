package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/models"
	"github.com/N4SJAMK/teamboard-api/utils"
)

func UpdateTicket(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex  = c.URLParams["board_id"]
		ticketIDHex = c.URLParams["ticket_id"]

		payload struct {
			Color   string `json:"color"`
			Heading string `json:"heading"`
			Content string `json:"content"`
		}
	)

	// read the payload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check that 'ticket_id' is a valid 'ObjectID'

	if !bson.IsObjectIdHex(ticketIDHex) {
		utils.Error(w, "'ticket_id' must be a valid 'ObjectID'",
			http.StatusBadRequest)
		return
	}
	ticketID := bson.ObjectIdHex(ticketIDHex)

	// get the board and make sure that the user is a member

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if !board.IsMember(user) {
		utils.Error(w, "", http.StatusForbidden)
		return
	}

	// apply the update and return the newly updated ticket

	var (
		ticket       = models.Ticket{}
		ticketQuery  = db.C("tickets").FindId(ticketID)
		ticketUpdate = mgo.Change{
			Update: bson.M{
				"color":   payload.Color,
				"heading": payload.Heading,
				"content": payload.Content,
			},
			ReturnNew: true,
		}
	)

	if _, err := ticketQuery.Apply(ticketUpdate, &ticket); err != nil {
		if err == mgo.ErrNotFound {
			utils.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, &ticket)
	return
}
