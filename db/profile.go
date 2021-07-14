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

func getProfiles(u User) (user User, err error) {
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
