package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var BackendURL string
var Token string

func init() {
	setUpBot()
}

func setUpBot() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	var exists bool

	BackendURL, exists = os.LookupEnv("BACKEND_URL")
	if !exists {
		log.Fatal("Backend URL not found")
	}

	Token, exists = os.LookupEnv("DISCORD_TOKEN")
	if !exists {
		log.Fatal("Token not found")
	}
}
