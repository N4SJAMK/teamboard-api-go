package models

import "time"
import "gopkg.in/mgo.v2/bson"

type (
	Board struct {
		ID bson.ObjectId `bson:"_id" json:"id"`

		Name        string `bson:"name"        json:"name"`
		Description string `bson:"description" json:"description"`

		Size       Size   `bson:"size"       json:"size"`
		Background string `bson:"background" json:"background"`

		Members []Member `bson:"members" json:"-"`

		CreatedAt time.Time     `bson:"created_at" json:"createdAt"`
		CreatedBy bson.ObjectId `bson:"created_by" json:"createdBy"`
	}

	Size struct {
		Width  int `bson:"width"  json:"width"`
		Height int `bson:"height" json:"height"`
	}

	Member struct {
		Role        string        `bson:"role"         json:"role"`
		UserID      bson.ObjectId `bson:"user_id"      json:"userID"`
		MemberSince time.Time     `bson:"member_since" json:"memberSince"`
	}
)

func (b *Board) IsMember(userID bson.ObjectId) bool {
	for _, member := range b.Members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

func (b *Board) GetRole(userID bson.ObjectId) string {
	for _, member := range b.Members {
		if member.UserID == userID {
			return member.Role
		}
	}
	return ""
}

func (b *Board) GetMember(userID bson.ObjectId) *Member {
	for _, member := range b.Members {
		if member.UserID == userID {
			return &member
		}
	}
	return nil
}
