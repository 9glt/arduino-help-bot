package main

import (
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

	// register defined functions in functions.go here
	fns.Bind("!help", fnHelp)
	fns.Bind("!tag", fnTag)

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
