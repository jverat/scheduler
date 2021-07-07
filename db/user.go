package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

//User is intended to keep the personalized configurations of schedules
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Password  string             `bson:"password" json:"password"`
	Profiles  Profiles           `bson:"profiles" json:"profiles"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Users []User

type Filter map[string]interface{}

func Create(u *User) (err error) {

	var searchResult User

	searchResult.Name = u.Name

	err = Read(&searchResult)

	if err != mongo.ErrNoDocuments {
		return err
	} else if searchResult.ID != primitive.NilObjectID {
		return fmt.Errorf("object already exists in db name: %s", searchResult.Name)
	}

	//Insertion
	_, err = usersCollection.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	//Return the created document, because the ObjectID is created by Mongo, the u.Name is what's going to be used to find it
	err = Read(u)

	return
}

func Read(u *User) (err error) {
	if u.ID == primitive.NilObjectID && u.Name == "" {
		return fmt.Errorf("id and name nil")
	}

	conditions := Filter{
		"$or": []Filter{
			{"_id": Filter{"$eq": u.ID}},
			{"name": Filter{"$eq": u.Name}},
		},
	}

	fmt.Printf("beforefind: %+v", u)
	userResult := usersCollection.FindOne(ctx, conditions)

	fmt.Printf("afterfind: %+v", userResult)

	err = userResult.Decode(&u)

	if err == mongo.ErrNoDocuments {
		err = fmt.Errorf("user %s not found", u.Name)
	}

	return
}

func Update(u *User) (err error) {
	result, err := usersCollection.ReplaceOne(
		ctx,
		Filter{"_id": Filter{"$eq": u.ID}},
		bson.M{
			"name":       u.Name,
			"password":   u.Password,
			"profiles":   u.Profiles,
			"updated_at": time.Now(),
		},
	)

	if result == nil {
		return fmt.Errorf("update result nil")
	}

	return Read(u)
}

func (u *User) Delete() (err error) {
	result, err := usersCollection.DeleteOne(
		ctx,
		Filter{"_id": Filter{"$eq": u.ID}},
	)
	if result == nil || result.DeletedCount != 0 {
		return fmt.Errorf("result = %v\n", result)
	}
	return
}
