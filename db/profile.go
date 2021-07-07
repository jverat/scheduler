package db

import (
	"time"
)

/*Profile is a set of data that configures the schedule stating how long will be every period of work (Workblocks), the rest that follows (Restblock),
the subsequent and supposedly longer break (LongRestBlock), and how many Workblocks before it happens*/
type Profile struct {
	//UserID                UserID    `bson:"_id"`
	Name                  string    `json:"name"`
	WorkblockDuration     int       `json:"workblock_duration"`
	RestblockDuration     int       `json:"restblock_duration"`
	LongRestblockDuration int       `json:"long_restblock_duration"`
	NWorkblocks           int       `json:"n_workblocks"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
}

type Profiles []Profile
