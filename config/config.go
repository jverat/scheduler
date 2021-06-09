package config

import (
	"github.com/joho/godotenv"
	"os"
)

var MongoDatabase string
var MongoUser string
var MongoPassword string

func SettingEnv() {

	godotenv.Load("./.env")

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
}
