package main

import "os"
import "github.com/zenazn/goji"

const (
	DEFAULT_DB_URL  = "mongodb://localhost"
	DEFAULT_DB_NAME = "teamboard-dev-go"
)

func main() {
	var (
		dburl  = os.Getenv("MONGODB_URL")
		dbname = os.Getenv("MONGODB_NAME")
	)

	if len(dburl) == 0 {
		dburl = DEFAULT_DB_URL
	}

	if len(dbname) == 0 {
		dbname = DEFAULT_DB_NAME
	}

	mdb := NewMongoDB(dburl, dbname)
	defer mdb.session.Close()

	goji.Get("/users",
		GetUsers)
	goji.Get("/users/:user_id",
		GetUser)

	goji.Get("/boards",
		GetBoards)
	goji.Get("/boards/:board_id",
		GetBoard)
	goji.Get("/boards/:board_id/members",
		GetMembers)
	goji.Get("/boards/:board_id/tickets",
		GetTickets)

	goji.Post("/boards",
		AddBoard)
	goji.Post("/boards/:board_id/members",
		AddMember)
	goji.Post("/boards/:board_id/tickets",
		AddTicket)

	goji.Delete("/boards/:board_id",
		RemoveBoard)
	goji.Delete("/boards/:board_id/members/:member_id",
		RemoveMember)
	goji.Delete("/boards/:board_id/tickets/:ticket_id",
		RemoveTicket)

	goji.Put("/boards/:board_id",
		UpdateBoard)
	goji.Put("/boards/:board_id/size",
		UpdateBoardSize)
	goji.Put("/boards/:board_id/members/:member_id",
		UpdateBoardMember)
	goji.Put("/boards/:board_id/tickets/:ticket_id",
		UpdateTicket)
	goji.Put("/boards/:board_id/tickets/:ticket_id/position",
		UpdateTicketPosition)

	goji.Use(mdb.Middleware)
	goji.Use(authenticate)

	goji.Serve()
	return
}
