package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*Profile is a set of data that configures the schedule stating how long will be every period of work (Workblocks), the rest that follows (Restblock),
the subsequent and supposedly longer break (LongRestBlock), and how many Workblocks before it happens*/
type Profile struct {
	UserID                UserID    `bson:"_id"`
	Name                  string    `json:"name"`
	WorkblockDuration     int       `json:"workblock_duration"`
	RestblockDuration     int       `json:"restblock_duration"`
	LongRestblockDuration int       `json:"long_restblock_duration"`
	NWorkblocks           int       `json:"n_workblocks"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
}

type ProfileID string

type Profiles []Profile

func (pf *Profile) Create() (uUpdateResult *mongo.UpdateResult, err error) {
	/*var u User
	userResult := usersCollection.FindOne(ctx, bson.D{{ Key: "id", Value: pf.UserID }})
	if userResult == nil {
		err = fmt.Errorf("User %s not found", pf.UserID)
		return
	}
	err = userResult.Decode(&u)
	if err != nil {
		return
	}

	fmt.Printf("User = %+v\n", u)

	u.Profiles = append(u.Profiles, pf)

	uUpdateResult, err = usersCollection.UpdateOne(ctx, bson.D{{ Key: "id", Value: pf.UserID }}, u)
	if err != nil {
		return
	}
	fmt.Printf("Profile inserted: %s\n", uUpdateResult)

	pfInsertionResult, err = profilesCollection.InsertOne(ctx, pf)
	if err != nil {
		return
	}
	fmt.Printf("Profile inserted: %s\n", pfInsertionResult)
	return*/

	u, err := pf.UserID.Read()
	if err != nil {
		return
	}

	for _, p := range u.Profiles {
		if p.Name == pf.Name {
			err = fmt.Errorf("Profile already exist")
			return
		}
	}

	u.Profiles = append(u.Profiles, *pf)

	return u.UpdateUser()
}

func (pfs Profiles) Read(userID string) (err error) {
	cursor, err := usersCollection.Find(ctx, bson.D{{Key: "userid", Value: userID}}, nil)
	if err != nil {
		return
	}

	err = cursor.All(ctx, &pfs)
	if err != nil {
		return
	}

	fmt.Println("Profiles retrieved")

	return
}

func (pf Profile) Update() (err error) {

	return
}

func (pf Profile) Delete() (err error) {

	return
}
