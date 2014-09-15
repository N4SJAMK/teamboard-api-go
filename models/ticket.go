package models

import "time"
import "gopkg.in/mgo.v2/bson"

type (
	Position struct {
		X int `bson:"x" json:"x"`
		Y int `bson:"y" json:"y"`
		Z int `bson:"z" json:"z"`
	}

	Ticket struct {
		ID        bson.ObjectId `bson:"_id"        json:"id"`
		BoardID   bson.ObjectId `bson:"board_id"   json:"-"`
		CreatedBy bson.ObjectId `bson:"created_by" json:"createdBy"`
		CreatedAt time.Time     `bson:"created_at" json:"createdAt"`

		Color   string `bson:"color"   json:"color"`
		Heading string `bson:"heading" json:"heading"`
		Content string `bson:"content" json:"content"`

		Position Position `bson:"position" json:"position"`
	}
)
