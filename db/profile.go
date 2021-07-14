package db

import "fmt"

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

func createProfiles(uID int, pf Profiles) (err error) {

	if len(pf) > 0 {
		for i := 0; i < len(pf); i++ {
			query := fmt.Sprintf("INSERT INTO public.profile (name, workblock_duration, restblock_duration, longrestblock_duration, n_workblocks, user_id) VALUES ('%s', %d, %d, %d, %d, %d)",
				pf[i].Name, pf[i].WorkblockDuration, pf[i].RestblockDuration, pf[i].LongRestblockDuration, pf[i].NWorkblocks, uID)
			_, err = Connection.Query(ctx, query)
			if err != nil {
				return
			}
		}
	}

	return
}

func readProfiles(u User) (user User, err error) {
	query := fmt.Sprintf("SELECT * FROM public.profile WHERE user_id = %d", u.ID)
	rows, err := Connection.Query(ctx, query)
	if err != nil {
		return
	}
	user = u
	for rows.Next() {
		var pf Profile
		err = rows.Scan(&pf.Name, &pf.WorkblockDuration, &pf.RestblockDuration, &pf.LongRestblockDuration, &pf.NWorkblocks, nil, nil)
		if err != nil {
			return
		}

		user.Profiles = append(user.Profiles, pf)
	}

	return
}
