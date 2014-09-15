package models

import "time"
import "gopkg.in/mgo.v2/bson"

type User struct {
	ID           bson.ObjectId `bson:"_id"           json:"id"`
	Email        string        `bson:"email"         json:"email"`
	Username     string        `bson:"username"      json:"username"`
	RegisteredAt time.Time     `bson:"registered_at" json:"registeredAt"`
}
