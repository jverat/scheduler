package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

type Filter map[string]interface{}

func (u *User) Create() (err error) {

	//Checks if the username is in use
	if sR := usersCollection.FindOne(ctx, bson.D{{Key: "name", Value: u.Name}}); sR != nil {
		return fmt.Errorf("%s already exist", u.Name)
	}

	//Insertion
	fmt.Printf("%+v\n", u)
	_, err = usersCollection.InsertOne(ctx, u)
	if err != nil {
		fmt.Println("err 1")
		return err
	}

	//Return the created document, because the ObjectID is created by Mongo. The u.Name is what's going to be used to find it
	err = u.Read()

	return
}

func (u *User) Read() (err error) {
	if primitive.ObjectID(u.ID) == primitive.NilObjectID && u.Name == "" {
		return fmt.Errorf("id and name nil")
	}

	conditions := Filter{
		"$or": []Filter{
			{"_id": Filter{"$eq": primitive.ObjectID(u.ID)}},
			{"name": Filter{"$eq": u.Name}},
		},
	}
	userResult := usersCollection.FindOne(ctx, conditions)
	if userResult == nil {
		err = fmt.Errorf("user %s not found", u.Name)
		return
	}

	err = userResult.Decode(&u)
	if err != nil {
		return
	}

	fmt.Printf("User = %+v\n", u)
	return
}

func (u *User) UpdateUser() (err error) {
	result, err := usersCollection.ReplaceOne(
		ctx,
		bson.M{"_id": Filter{"$eq": primitive.ObjectID(u.ID)}},
		bson.M{
			"Name":      u.Name,
			"Password":  u.Password,
			"Profiles":  u.Profiles,
			"UpdatedAt": time.Now(),
		},
	)

	if result == nil {
		return fmt.Errorf("update result nil")
	}

	return u.Read()
}

func (u *User) DeleteUser() (err error) {
	result, err := usersCollection.DeleteOne(
		ctx,
		bson.M{"_id": Filter{"$eq": primitive.ObjectID(u.ID)}},
	)
	if result == nil || result.DeletedCount != 0 {
		return fmt.Errorf("result = %v\n", result)
	}
	return
}
