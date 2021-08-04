package main

import (
	"database/sql"
	"github.com/jwambugu/go-simple-bank-class/api"
	db "github.com/jwambugu/go-simple-bank-class/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:3000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to the db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	log.Fatal(server.Start(serverAddress))
}
