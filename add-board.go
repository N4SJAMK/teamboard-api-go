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

func AddBoard(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db   = c.Env["db"].(*mgo.Database)
		user = c.Env["user"].(bson.ObjectId)

		payload struct {
			Name        string      `json:"name"`
			Description string      `json:"description"`
			Size        models.Size `json:"size"`
			Background  string      `json:"background"`
		}
	)

	// read the payload

	if err := utils.ReadJSON(r, &payload); err != nil {
		utils.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create the board, with the requesting user as the creator

	board := models.Board{
		ID: bson.NewObjectId(),

		CreatedBy: user,
		CreatedAt: time.Now(),

		Name:        payload.Name,
		Size:        payload.Size,
		Background:  payload.Background,
		Description: payload.Description,

		Members: []models.Member{
			{
				Role:        "admin",
				UserID:      user,
				MemberSince: time.Now(),
			},
		},
	}

	if err := db.C("boards").Insert(&board); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// respond with the created board

	utils.WriteJSON(w, http.StatusCreated, &board)
	return
}
