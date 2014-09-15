package main

import (
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/zenazn/goji/web"
)

type MongoDB struct {
	name    string
	session *mgo.Session
}

func (mdb *MongoDB) Middleware(c *web.C, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := mdb.session.Copy()
		defer session.Close()

		c.Env["db"] = session.DB(mdb.name)

		h.ServeHTTP(w, r)
		return
	})
}

func NewMongoDB(url, name string) *MongoDB {
	session, err := mgo.Dial(url)

	if err != nil {
		panic(err)
	}

	return &MongoDB{
		name:    name,
		session: session,
	}
}
