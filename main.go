package main

import (
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	envToken   = os.Getenv("BOT_TOKEN")
	envRoles   = os.Getenv("BOT_ADMIN_ROLES")
	envDocsDir = os.Getenv("BOT_DOCS_DIR")

	roles = make(map[string]struct{})

	fns *Functions

	tagsRegistry     = make(map[string]*Tag)
	tagsRegistryLock = &sync.RWMutex{}
)

func main() {
	if envToken == "" {
		panic("BOT_TOKEN is not set")
	}

	if envDocsDir == "" {
		envDocsDir = "/docs"
	}

	for _, role := range strings.Split(envRoles, ",") {
		roles[role] = struct{}{}
	}

	ScanForTags()

	fns = NewFunctions(10)

	// register defined functions in functions.go here
	fns.Bind("!help", fnHelp)
	fns.Bind("!tag", fnTag)
	fns.Bind("!reload", fnReload)

	fns.Fallback(fnFallback)

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

	if strings.HasPrefix(m.Content, "!") && checkLen(m.Content) {
		fns.Run(m.Content, s, m)
	}

}
