package test

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"scheduler/db"
	"testing"
	"time"
)

func TestUser_Create(t *testing.T) {

	err := db.DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.DatabaseTurnoff()
		if err != nil {
			log.Fatal(err)
		}
	}()

	user := db.User{
		Name:      "fgv",
		Password:  "",
		Profiles:  db.Profiles{},
		CreatedAt: time.Now(),
	}
	err = db.Create(&user)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
		t.Fail()
	} else if user.ID == primitive.NilObjectID {
		t.Errorf("uID doesn't exist")
		t.Fail()
	} else {
		fmt.Printf("%+v\n", user)
		return
	}

}

func TestUser_CreateWithProfile(t *testing.T) {
	err := db.DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = db.DatabaseTurnoff()
		if err != nil {
			log.Fatal(err)
		}
	}()

	user := db.User{
		Name:     "userWithProfile",
		Password: "",
		Profiles: db.Profiles{db.Profile{
			Name:                  "",
			WorkblockDuration:     0,
			RestblockDuration:     0,
			LongRestblockDuration: 0,
			NWorkblocks:           0,
			CreatedAt:             time.Now(),
		}},
		CreatedAt: time.Now(),
	}
	err = db.Create(&user)

	if err != nil {
		fmt.Printf("username: %s\n", user.Name)
		t.Errorf("Error: %s", err.Error())
		t.Fail()
	} else if user.ID == primitive.NilObjectID {
		t.Errorf("uID doesn't exist")
		t.Fail()
	} else {
		return
	}
}

func TestUser_Read(t *testing.T) {

	userID, err := primitive.ObjectIDFromHex("60e4cce3351f0bf2ec9edd0f")
	fmt.Println(userID)

	if err != nil {
		t.Errorf("??????? %s", err)
		t.Fail()
	}

	myUser := db.User{
		ID:       primitive.NilObjectID,
		Name:     "ReadTestUser",
		Password: "",
		Profiles: db.Profiles{},
	}

	fmt.Printf("%+v\n", myUser)

	err = db.Read(&myUser)

	if err != nil {
		t.Errorf("Error: %s", err)
		t.Fail()
	} else if myUser.ID == primitive.NilObjectID {
		t.Errorf("Error: user not found")
		t.Fail()
	}

	myUser.ID = userID
	myUser.Name = ""

	fmt.Printf("%+v\n", myUser)

	err = db.Read(&myUser)

	fmt.Printf("%+v\n", myUser)

	if err != nil {
		t.Errorf("Error: %s", err)
		t.Fail()
	} else if myUser.ID != userID {
		t.Errorf("Weird uID")
		t.Fail()
	}

	return
}
