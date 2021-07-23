package main

import (
	"log"

	"scheduler/db"
)

/*
TODO
DB error handling
REST routes
Learn algorithmics
*/
func main() {
	err := db.DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db.Connection.Close()
		db.CancelFunc()
		return
	}()
}
