package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var PostgresDatabase string
var PostgresUser string
var PostgresPassword string
var PostgresHost string

func SettingEnv() {

	absPath, err := os.Getwd()
	var dotenv string
	if strings.HasSuffix(absPath, "scheduler") {
		dotenv = ".env"
	} else {
		dotenv = "../.env"
	}

	err = godotenv.Load(dotenv)
	if err != nil {
		fmt.Println(err)
	}

	var exists bool

	PostgresDatabase, exists = os.LookupEnv("POSTGRES_DATABASE")
	if !exists {
		PostgresDatabase = "scheduler-db"
	}

	PostgresUser, exists = os.LookupEnv("POSTGRES_USER")
	if !exists {
		PostgresUser = "admin"
	}

	PostgresPassword, exists = os.LookupEnv("POSTGRES_PASSWORD")
	if !exists {
		PostgresPassword = "admin"
	}

	PostgresHost, exists = os.LookupEnv("POSTGRES_HOST")
	if !exists {
		PostgresHost = "localhost"
	}
}
