package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	envToken = os.Getenv("BOT_TOKEN")

	fns *Functions
)

func main() {
	if envToken == "" {
		panic("BOT_TOKEN is not set")
	}
	fns = NewFunctions(10)
	fns.Bind("!help", func(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
		_, err := ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Hello World! %s", s))
		if err != nil {
			log.Printf("%v", err)
		}
	})
	dg, err := discordgo.New("Bot " + envToken)
	if err != nil {
		panic(err)
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentGuildMessages
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	log.Printf("Up and Running!")

	runtime.Goexit()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		log.Printf("Spot tagged message: %v", m.Content)
		fns.Run(m.Content, s, m)
	}

}

func update() {

}
