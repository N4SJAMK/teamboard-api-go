package models

import "gopkg.in/mgo.v2/bson"

type Token struct {
	Secret string        `bson:"secret"`
	UserID bson.ObjectId `bson:"user_id"`
}
