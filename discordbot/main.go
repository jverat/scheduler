package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"scheduler/discordbot/config"
)

var Session *discordgo.Session
var BotID string

func main() {
	config.SetUpBot()
	var err error
	Session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
	}

	u, err := Session.User("@me")
	if err != nil {
		log.Printf("%s\n", err)
	}

	BotID = u.ID
	err = Session.Open()
	if err != nil {
		log.Fatal("error opening connection: ", err)
	}
	log.Printf("Bot running!")
	<-make(chan int)
}
