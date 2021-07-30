package db

import (
	"fmt"
	"github.com/jackc/pgx/v4"
)

/*Profile is a set of data that configures the schedule stating how long will be every period of work (Workblocks), the rest that follows (Restblock),
the subsequent and supposedly longer break (LongRestBlock), and how many Workblocks before it happens*/
type Profile struct {
	ID                    int
	Name                  string
	WorkblockDuration     int
	RestblockDuration     int
	LongRestblockDuration int
	NWorkblocks           int
}

type Profiles []Profile

func (profiles Profiles) contains(profile Profile) bool {
	for _, pf := range profiles {
		if pf.ID == profile.ID {
			return true
		}
	}
	return false
}

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

func orderByID(input Profiles) (result Profiles) {
	if len(input) == 0 {
		return input
	}

	var minorIndex, minorID int

	minorID = input[0].ID

	for i, p := range input {
		if p.ID < minorID {
			minorID = p.ID
			minorIndex = i
		}
	}

	result = append(result, input[minorIndex])

	for len(result) < len(input) {
		minorID++
		for _, pf := range input {
			if pf.ID == minorID {
				result = append(result, pf)
			}
		}
	}
	return
}

func CreateProfiles(uID int, pf Profiles) (profiles Profiles, err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)

	if len(pf) > 0 {
		for i := 0; i < len(pf); i++ {
			query := fmt.Sprintf("INSERT INTO public.\"profile\" (name, workblock_duration, restblock_duration, longrestblock_duration, n_workblocks, user_id) VALUES ('%s', %d, %d, %d, %d, %d)",
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

	return ReadProfiles(uID)
}

func ReadProfiles(uID int) (profiles Profiles, err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)
	query := fmt.Sprintf("SELECT * FROM public.\"profile\" WHERE user_id = %d", uID)
	queryChan <- query
	close(queryChan)

	rows := <-outputChan

	for rows.Next() {
		var pf Profile
		err = rows.Scan(&pf.Name, &pf.WorkblockDuration, &pf.RestblockDuration, &pf.LongRestblockDuration, &pf.NWorkblocks, &pf.ID, nil)
		if err != nil {
			return
		}

		profiles = append(profiles, pf)
	}

	profiles = orderByID(profiles)
	return
}

func updateProfiles(uID int, newProfiles Profiles, oldProfiles Profiles) (err error) {
	queryUpdate := "UPDATE public.\"profile\" SET name = '%s', workblock_duration = %d, restblock_duration = %d, longrestblock_duration = %d, n_workblocks = %d WHERE id = %d"
	queryCreate := "INSERT INTO public.\"profile\" (name, workblock_duration, restblock_duration, longrestblock_duration, n_workblocks, id, user_id) VALUES ('%s', %d, %d, %d, %d, DEFAULT, %d) RETURNING id"

	queryChan, outputChan, errChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errChan)

	for i, pf := range newProfiles {
		if pf.ID == 0 {
			queryChan <- fmt.Sprintf(queryCreate, pf.Name, pf.WorkblockDuration, pf.RestblockDuration, pf.LongRestblockDuration, pf.NWorkblocks, uID)
			select {
			case err = <-errChan:
				return
			case r := <-outputChan:
				for r.Next() {
					err = r.Scan(&newProfiles[i].ID)
				}
				if err != nil {
					return
				}
			}
		}
	}

	oldProfiles, err = ReadProfiles(uID)
	if err != nil {
		return
	}

	newProfiles = orderByID(newProfiles)
	oldProfiles = orderByID(oldProfiles)

	if len(newProfiles) == len(oldProfiles) {
		for i := range newProfiles {
			if newProfiles[i].ID == oldProfiles[i].ID && !identical(newProfiles[i], oldProfiles[i]) {
				q := fmt.Sprintf(queryUpdate, newProfiles[i].Name, newProfiles[i].WorkblockDuration, newProfiles[i].RestblockDuration, newProfiles[i].LongRestblockDuration, newProfiles[i].NWorkblocks, newProfiles[i].ID)
				queryChan <- q
				select {
				case err = <-errChan:
					return
				case _ = <-outputChan:
				}
			}
		}
		close(queryChan)
	} else if len(newProfiles) < len(oldProfiles) {
		oldProfiles, err = deleteFromUpdateProfiles(oldProfiles, newProfiles)
		if err != nil {
			return
		}

		if !identicals(oldProfiles, newProfiles) {
			return updateProfiles(uID, newProfiles, oldProfiles)
		}
	}
	return
}

func deleteFromUpdateProfiles(oldProfiles, newProfiles Profiles) (profiles Profiles, err error) {
	for i := range newProfiles {
		if newProfiles[i].ID != oldProfiles[i].ID {
			err = DeleteProfile(oldProfiles[i])
			if err != nil {
				return oldProfiles, fmt.Errorf("error while deleting profiles: %s", err)
			}
			if len(newProfiles) != len(oldProfiles) {
				return deleteFromUpdateProfiles(append(oldProfiles[:i], oldProfiles[i+1:]...), newProfiles)
			} else {
				return append(oldProfiles[:i], oldProfiles[i+1:]...), err
			}
		}
	}

	return
}

func DeleteProfile(pf Profile) (err error) {
	queryChan, outputChan, errorChan := make(chan string), make(chan pgx.Rows), make(chan error)
	go AcquireConn(queryChan, outputChan, errorChan)
	query := fmt.Sprintf("DELETE FROM public.\"profile\" * WHERE id = '%d'", pf.ID)
	queryChan <- query
	close(queryChan)
	select {
	case err = <-errorChan:
		return
	case _ = <-outputChan:
		return
	}
}

func DeleteProfiles(uID int) (err error) {
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
