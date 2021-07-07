package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var MongoDatabase string
var MongoUser string
var MongoPassword string
var MongoHost string

func SettingEnv() {

	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println(err)
	}

	var exists bool

	MongoDatabase, exists = os.LookupEnv("MONGO_DATABASE")
	if !exists {
		MongoDatabase = "scheduler-db"
	}

	MongoUser, exists = os.LookupEnv("MONGO_USER")
	if !exists {
		MongoUser = "admin"
	}

	MongoPassword, exists = os.LookupEnv("MONGO_PASSWORD")
	if !exists {
		MongoPassword = "admin"
	}

	MongoHost, exists = os.LookupEnv("MONGO_HOST")
	if !exists {
		MongoHost = "localhost"
	}
}
