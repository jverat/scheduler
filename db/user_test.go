package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
	"time"
)

func TestUser_Create(t *testing.T) {

	err := DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = DatabaseTurnoff()
		if err != nil {
			log.Fatal(err)
		}
	}()

	user := User{
		Name:      "fgv",
		Password:  "",
		Profiles:  Profiles{},
		CreatedAt: time.Now(),
	}
	err = Create(&user)

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
	err := DatabaseSetUp()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = DatabaseTurnoff()
		if err != nil {
			log.Fatal(err)
		}
	}()

	user := User{
		Name:     "userWithProfile",
		Password: "",
		Profiles: Profiles{Profile{
			Name:                  "",
			WorkblockDuration:     0,
			RestblockDuration:     0,
			LongRestblockDuration: 0,
			NWorkblocks:           0,
			CreatedAt:             time.Now(),
		}},
		CreatedAt: time.Now(),
	}
	err = Create(&user)

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

	myUser := User{
		ID:       primitive.NilObjectID,
		Name:     "ReadTestUser",
		Password: "",
		Profiles: Profiles{},
	}

	fmt.Printf("%+v\n", myUser)

	err = Read(&myUser)

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

	err = Read(&myUser)

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
