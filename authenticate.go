package main

import (
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"

	"github.com/N4SJAMK/teamboard-api/models"
	"github.com/N4SJAMK/teamboard-api/utils"
)

func authenticate(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			db   = c.Env["db"].(*mgo.Database)
			auth = strings.Split(r.Header.Get("authorization"), " ")
		)

		// validate that the authorization header is correctly formatted

		if len(auth) != 2 || auth[0] != "Bearer" {
			utils.Error(w, "Authorization: Bearer <Token>",
				http.StatusUnauthorized)
			return
		}

		// find the token

		var (
			token      = models.Token{}
			tokenQuery = db.C("tokens").Find(bson.M{
				"secret": auth[1],
			})
		)

		if err := tokenQuery.One(&token); err != nil {
			if err == mgo.ErrNotFound {
				utils.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			utils.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// find the user corresponding to the token and attach its
		// 'ID' attribute to the request's context under the 'user' key

		var (
			user      = models.User{}
			userQuery = db.C("users").FindId(token.UserID)
		)

		if err := userQuery.One(&user); err != nil {
			if err == mgo.ErrNotFound {
				utils.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			utils.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Env["user"] = user.ID

		h.ServeHTTP(w, r)
		return
	})
}
