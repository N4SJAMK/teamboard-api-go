package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func RemoveTicket(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex  = c.URLParams["board_id"]
		ticketIDHex = c.URLParams["ticket_id"]
	)

	// get the board for this request, and make sure the user making the
	// request is a member on this board

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if !board.IsMember(user) {
		utils.Error(w, "", http.StatusForbidden)
		return
	}

	// get the ticket

	ticket, err := getTicket(db, boardIDHex, ticketIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	// remove the ticket from 'tickets'-collection and return the removed ticket

	if err := db.C("tickets").RemoveId(ticket.ID); err != nil {
		if err == mgo.ErrNotFound {
			utils.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, ticket)
	return
}
