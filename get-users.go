package main

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/models"
	"github.com/N4SJAMK/teamboard-api/utils"
)

func GetUsers(c web.C, w http.ResponseWriter, r *http.Request) {
	var (
		db         = c.Env["db"].(*mgo.Database)
		users      = []models.User{}
		usersQuery = db.C("users").Find(nil)
	)

	if err := usersQuery.All(&users); err != nil {
		utils.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
	return
}
