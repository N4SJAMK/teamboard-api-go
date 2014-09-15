package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/N4SJAMK/teamboard-api/models"
)

type APIError struct {
	message        string
	HTTPStatusCode int
}

func (err *APIError) Error() string {
	if len(err.message) > 0 {
		return err.message
	}
	return http.StatusText(err.HTTPStatusCode)
}

func getBoard(db *mgo.Database, boardIDHex string) (*models.Board, *APIError) {
	if !bson.IsObjectIdHex(boardIDHex) {
		return nil, &APIError{"'board_id' must be a valid 'ObjectID'",
			http.StatusBadRequest}
	}
	boardID := bson.ObjectIdHex(boardIDHex)

	var (
		board      = models.Board{}
		boardQuery = db.C("boards").FindId(boardID)
	)

	if err := boardQuery.One(&board); err != nil {
		if err == mgo.ErrNotFound {
			return nil, &APIError{"Board not found", http.StatusNotFound}
		}
		return nil, &APIError{err.Error(), http.StatusInternalServerError}
	}

	return &board, nil
}

func getUser(db *mgo.Database, userIDHex string) (*models.User, *APIError) {
	if !bson.IsObjectIdHex(userIDHex) {
		return nil, &APIError{"'user_id' must be a valid 'ObjectID'",
			http.StatusBadRequest}
	}
	userID := bson.ObjectIdHex(userIDHex)

	var (
		user      = models.User{}
		userQuery = db.C("users").FindId(userID)
	)

	if err := userQuery.One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, &APIError{"User not found", http.StatusNotFound}
		}
		return nil, &APIError{err.Error(), http.StatusInternalServerError}
	}

	return &user, nil
}

func getTicket(db *mgo.Database, boardIDHex, ticketIDHex string) (
	*models.Ticket, *APIError,
) {
	if !bson.IsObjectIdHex(boardIDHex) || !bson.IsObjectIdHex(ticketIDHex) {
		return nil, &APIError{
			"Both 'board_id' and 'ticket_id' must be a valid 'ObjectIDs'",
			http.StatusBadRequest}
	}

	var (
		ticket      = models.Ticket{}
		ticketQuery = db.C("tickets").Find(bson.M{
			"_id":      bson.ObjectIdHex(ticketIDHex),
			"board_id": bson.ObjectIdHex(boardIDHex),
		})
	)

	if err := ticketQuery.One(&ticket); err != nil {
		if err == mgo.ErrNotFound {
			return nil, &APIError{"Ticket not found", http.StatusNotFound}
		}
		return nil, &APIError{err.Error(), http.StatusInternalServerError}
	}

	return &ticket, nil
}

func updateBoard(
	db *mgo.Database, userID bson.ObjectId, boardIDHex string, change bson.M,
) (
	*models.Board,
	*APIError,
) {
	// get the board
	// check that the user is a member, and has admin access

	board, err := getBoard(db, boardIDHex)
	if err != nil {
		return nil, err
	}

	if board.GetRole(userID) != "admin" {
		return nil, &APIError{"User must have 'admin' privileges",
			http.StatusForbidden}
	}

	// apply updates and return the newly updated board

	var (
		boardQuery  = db.C("boards").FindId(board.ID)
		boardUpdate = mgo.Change{
			Update:    change,
			ReturnNew: true,
		}
	)

	if _, err := boardQuery.Apply(boardUpdate, board); err != nil {
		if err == mgo.ErrNotFound {
			return nil, &APIError{"Board not found", http.StatusNotFound}
		}
		return nil, &APIError{err.Error(), http.StatusInternalServerError}
	}

	return board, nil
}
