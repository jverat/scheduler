package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var Token string

func SetUpBot() {

	absPath, _ := os.Getwd()
	var dotenv string
	if strings.HasSuffix(absPath, "scheduler/discordbot") {
		dotenv = ".env"
	} else {
		dotenv = "../.env"
	}
	log.Printf("%s\n", dotenv)

	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	var exists bool
	Token, exists = os.LookupEnv("DISCORD_TOKEN")
	if !exists {
		log.Fatal("Token not found")
	}
}
