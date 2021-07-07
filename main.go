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
		err = db.DatabaseTurnoff()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
