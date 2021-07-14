package main

import (
	"fmt"
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

	u, err := db.Create(db.User{Name: "skere2", Password: "skere otra vez",
		Profiles: []db.Profile{{
			Name:                  "unperfil",
			WorkblockDuration:     6,
			RestblockDuration:     6,
			LongRestblockDuration: 6,
			NWorkblocks:           6,
		}, {
			Name:                  "otroperfil",
			WorkblockDuration:     7,
			RestblockDuration:     7,
			LongRestblockDuration: 7,
			NWorkblocks:           7,
		}}})
	fmt.Printf("%+v\nErr=%s", u, err)
}
