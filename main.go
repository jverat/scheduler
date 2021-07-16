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
	}()

	/*u, err := db.Create(db.User{ID: 6, Name: "skere2", Password: "skere otra vez",
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
	}, {
		Name:                  "untercerperfil",
		WorkblockDuration:     8,
		RestblockDuration:     8,
		LongRestblockDuration: 8,
		NWorkblocks:           8,
	}}})*/
	u := db.User{ID: 19}
	//err = db.Connection.QueryRow(context.Background(), "SELECT * FROM public.\"user\" WHERE id = 19").Scan(&u.ID, &u.Name, &u.Password)
	u, err = db.Read(u)
	fmt.Printf("%+v\nErr=%s\n", u, err)
}
