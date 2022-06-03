package main

import (
	"database/sql"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error encountered: %s \\n", err)
		os.Exit(1)
	}
}

func run() error {

	// setup database
	// database is injected to storage as a dependency
	db, err := setUpDatabase(connectionString)

	// setup up storage
	// services depend on storage
	storage := repository.NewStorage(db)

	// setup services
	// this requires storage, interface between
	// service specific storage and server

	// run the migrations
	//

	// create the server. handles client requests and directs to correct
	// api endpoint

	// run the server

}

type User struct {
}

func setUpDatabase(connectionString string) (db *sql.DB, err error) {
	db = struct {
		Users map[int]User
	}{}
}
