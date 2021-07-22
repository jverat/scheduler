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
		return
	}()

	fmt.Println(db.Delete(db.User{ID: 44}))

	/*err = db.Update(db.User{ID: 44, Name: "skere3", Password: "skere otra vez",
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
		}}})

	fmt.Printf("Err: %s\n", err)*/

	/*pf := []db.Profile{{
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
	}}


	for i, p := range pf {
		conn, err := db.Connection.Acquire(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}

		query := fmt.Sprintf("INSERT INTO public.profile (name, workblock_duration, restblock_duration, longrestblock_duration, n_workblocks, user_id) VALUES ('%s', %d, %d, %d, %d, %d)",
			p.Name, p.WorkblockDuration, p.RestblockDuration, p.LongRestblockDuration, p.NWorkblocks, 44)

		_, err = conn.Query(context.Background(), query)
		if err != nil {
			fmt.Println("#", i, " ", err)
			return
		}
		conn.Release()
	}*/
}
