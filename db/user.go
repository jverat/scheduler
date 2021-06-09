package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

//User is intended to keep the personalized configurations of schedules
type User struct {
	ID        UserID    `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Profiles  Profiles  `json:"profiles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UserID primitive.ObjectID

type Users []User

func (u *User) Create() (err error) {
	fmt.Printf("%+v\n", u)
	result, err := usersCollection.InsertOne(ctx, u)
	if err != nil {
		fmt.Println("err 1")
		return err
	}

	fmt.Println("2")
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		var oidFromHex primitive.ObjectID
		oidFromHex, err = primitive.ObjectIDFromHex(oid.Hex())
		*u, err = UserID(oidFromHex).Read()
		return err
	} else {
		fmt.Println("err 2")
		return fmt.Errorf("no objectID")
	}
}

func (userID UserID) Read() (user User, err error) {
	userResult := usersCollection.FindOne(ctx, bson.D{{Key: "_id", Value: userID}})
	if userResult == nil {
		err = fmt.Errorf("user %s not found", userID)
		return
	}

	err = userResult.Decode(&user)
	if err != nil {
		return
	}

	fmt.Printf("User = %+v\n", user)
	return
}

func (u User) UpdateUser() (uUpdateResult *mongo.UpdateResult, err error) {
	return
}
