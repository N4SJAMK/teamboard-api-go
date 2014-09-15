package main

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/utils"
)

func GetUser(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db        = c.Env["db"].(*mgo.Database)
		userIDHex = c.URLParams["user_id"]
	)

	user, err := getUser(db, userIDHex)
	if err != nil {
		utils.Error(w, err.Error(), err.HTTPStatusCode)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
	return
}
