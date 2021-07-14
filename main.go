package main

import (
	"log"
	"scheduler/db"
)

func main() {
	err := db.DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db.Connection.Close()
		db.CancelFunc()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
