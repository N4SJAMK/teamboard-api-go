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

func AddTicket(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		boardIDHex = c.URLParams["board_id"]

		payload struct {
			Color    string          `json:"color"`
			Heading  string          `json:"heading"`
			Content  string          `json:"content"`
			Position models.Position `json:"position"`
		}
	)

	// read the payload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the required board and make sure the requesting user is a member

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	if !board.IsMember(user) {
		utils.Error(w, "User must be a board member", http.StatusForbidden)
		return
	}

	// insert a new 'Ticket' into 'tickets'

	newTicket := models.Ticket{
		ID:      bson.NewObjectId(),
		BoardID: board.ID,

		Color:   payload.Color,
		Heading: payload.Heading,
		Content: payload.Content,

		Position: payload.Position,

		CreatedBy: user,
		CreatedAt: time.Now(),
	}

	if err := db.C("tickets").Insert(&newTicket); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the added ticket

	utils.WriteJSON(w, http.StatusCreated, &newTicket)
	return
}
