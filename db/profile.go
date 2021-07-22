package db

import (
	"fmt"
	"github.com/jackc/pgx/v4"
)

/*Profile is a set of data that configures the schedule stating how long will be every period of work (Workblocks), the rest that follows (Restblock),
the subsequent and supposedly longer break (LongRestBlock), and how many Workblocks before it happens*/
type Profile struct {
	//UserID                UserID    `bson:"_id"`
	Name                  string `json:"Name"`
	WorkblockDuration     int    `json:"WorkblockDuration"`
	RestblockDuration     int    `json:"RestblockDuration"`
	LongRestblockDuration int    `json:"LongRestblockDuration"`
	NWorkblocks           int    `json:"NWorkblocks"`
}

type Profiles []Profile

func identical(x, y Profile) bool {
	if x.Name == y.Name && x.WorkblockDuration == y.WorkblockDuration && x.RestblockDuration == y.RestblockDuration && x.LongRestblockDuration == y.LongRestblockDuration && x.NWorkblocks == y.NWorkblocks {
		return true
	} else {
		return false
	}
}

func identicals(x, y Profiles) bool {
	if len(x) != len(y) {
		return false
	}

	for i := 0; i < len(x); i++ {
		if !identical(x[i], y[i]) {
			return false
		}
	}

	return true
}

func createProfiles(uID int, pf Profiles) (err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)

	if len(pf) > 0 {
		for i := 0; i < len(pf); i++ {
			query := fmt.Sprintf("INSERT INTO public.profile (name, workblock_duration, restblock_duration, longrestblock_duration, n_workblocks, user_id) VALUES ('%s', %d, %d, %d, %d, %d)",
				pf[i].Name, pf[i].WorkblockDuration, pf[i].RestblockDuration, pf[i].LongRestblockDuration, pf[i].NWorkblocks, uID)

			queryChan <- query
			select {
			case err = <-errorChan:
				{
					close(queryChan)
					return
				}
			case _ = <-outputChan:
			}
		}
	}
	close(queryChan)

	return
}

func readProfiles(uID int) (profiles Profiles, err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)
	query := fmt.Sprintf("SELECT * FROM public.\"profile\" WHERE user_id = %d", uID)
	queryChan <- query
	close(queryChan)

	rows := <-outputChan

	for rows.Next() {
		var pf Profile
		err = rows.Scan(&pf.Name, &pf.WorkblockDuration, &pf.RestblockDuration, &pf.LongRestblockDuration, &pf.NWorkblocks, nil, nil)
		if err != nil {
			return
		}

		profiles = append(profiles, pf)
	}

	return
}

func updateProfiles(uID int, newProfiles Profiles, oldProfiles Profiles) (err error) {
	query := "UPDATE public.profile SET name = '%s', workblock_duration = %d, restblock_duration = %d, longrestblock_duration = %d, n_workblocks = %d WHERE name = '%s' AND user_id = %d"

	queryChan, outputChan, errChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errChan)

	if len(newProfiles) == len(oldProfiles) {
		for i := 0; i < len(newProfiles); i++ {
			if !identical(newProfiles[i], oldProfiles[i]) {
				q := fmt.Sprintf(query, newProfiles[i].Name, newProfiles[i].WorkblockDuration, newProfiles[i].RestblockDuration, newProfiles[i].LongRestblockDuration, newProfiles[i].NWorkblocks, oldProfiles[i].Name, uID)
				queryChan <- q
				select {
				case err = <-errChan:
					return
				case _ = <-outputChan:
				}
			}
			close(queryChan)
		}
	} else if len(newProfiles) > len(oldProfiles) {
		var pf Profiles

		for i := 0; i < len(pf); i++ {
			if i >= len(oldProfiles) {

			}
		}
	}

	//conn.Release()

	return
}

func DeleteProfile(uID int, pf Profile) (err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)
	query := fmt.Sprintf("DELETE FROM public.\"profile\" * WHERE name = '%s' AND user_id = %d", pf.Name, uID)
	queryChan <- query
	close(queryChan)
	select {
	case err = <-errorChan:
		return
	case _ = <-outputChan:
		return
	}
}

func deleteProfiles(uID int) (err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)
	query := fmt.Sprintf("DELETE FROM public.\"profile\" * WHERE user_id = %d", uID)
	queryChan <- query
	close(queryChan)
	select {
	case err = <-errorChan:
		return
	case _ = <-outputChan:
		return
	}
}
