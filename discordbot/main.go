package main

import (
	"github.com/lus/dgc"
	"log"
	"os"
	"os/signal"
	"scheduler/discordbot/config"
	"scheduler/discordbot/querying"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Session *discordgo.Session
var BotID string

func init() {
	var err error
	Session, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
	}
}

var (
	rateLimiter = dgc.NewRateLimiter(5*time.Second, 1*time.Second, func(ctx *dgc.Ctx) {
		ctx.RespondText("chill bro")
	})
	commands = []*dgc.Command{
		{
			Name:        "sign-up",
			Description: "sign-up in scheduler! Create your account and start working with your own profiles",
			RateLimiter: rateLimiter,
			Handler: func(ctx *dgc.Ctx) {
				if ctx.Event.Author.ID == BotID {
					return
				}

				channel, err := ctx.Session.UserChannelCreate(ctx.Event.Author.ID)
				if err != nil {
					log.Printf("error creating channel: %s\n", err)
					_, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "Something went wrong while sending the DM!")
					if err != nil {
						log.Printf("error while chatting in server: %s\n", err)
					}
					return
				}
				_, err = ctx.Session.ChannelMessageSend(channel.ID, "lmao")
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"sign-up": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.User.ID == s.State.User.ID {
				return
			}

			channel, err := s.UserChannelCreate(i.User.ID)
			if err != nil {
				log.Printf("error creating channel: %s\n", err)
				_, err := s.ChannelMessageSend(i.ChannelID, "Something went wrong while sending the DM!")
				if err != nil {
					log.Printf("?: %s\n", err)
				}
				return
			}
			_, err = s.ChannelMessageSend(channel.ID, "lmao")
		},
	}
)

func init() {
	Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("skere")
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	log.Println(querying.Client)
	defer func() {
		err := Session.Close()
		if err != nil {
			log.Fatal("error closing session")
		}
	}()
	u, err := Session.User("@me")
	if err != nil {
		log.Printf("%s\n", err)
	}

	BotID = u.ID

	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	router := dgc.Create(&dgc.Router{
		Prefixes:         []string{"sch!"},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		PingHandler: func(ctx *dgc.Ctx) {
			ctx.RespondText("watchu want")
		},
		Commands: commands,
	})
	router.RegisterDefaultHelpCommand(Session, rateLimiter)
	err = Session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}
	router.Initialize(Session)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}
