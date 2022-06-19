package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	envToken         = os.Getenv("BOT_TOKEN")
	envRoles         = os.Getenv("BOT_ADMIN_ROLES")
	envDocsDir       = os.Getenv("BOT_DOCS_DIR")
	envBlacklistExts = os.Getenv("BOT_BLACKLIST_EXTS")

	varBlacklistExts = []string{}

	roles = make(map[string]struct{})

	fns *Functions

	tagsRegistry     = make(map[string]*Tag)
	tagsRegistryLock = &sync.RWMutex{}
	scannerLock      = NewLocker()
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

	if envBlacklistExts != "" {
		varBlacklistExts = strings.Split(envBlacklistExts, ",")
	}

	ScanForTags()

	fns = NewFunctions(10)

	// register defined functions in functions.go here
	fns.Bind("!help", fnHelp)
	fns.Bind("!tag", fnTag)
	fns.Bind("!reload", fnReload)
	fns.Bind("!ping", fnPing)

	fns.Fallback(fnFallback)

	dg, err := discordgo.New("Bot " + envToken)
	if err != nil {
		panic(err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("mh")
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	dg.Identify.Presence.Game.Name = "!help | no version"
	dg.Identify.Presence.Game.Type = 3
	dg.Identify.Intents = discordgo.IntentGuildMessages

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "275273435930951680", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
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

	for _, attachment := range m.Attachments {
		if !checkExt(attachment.Filename) {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("disallowed filetype %v", attachment.Filename))
			if err != nil {
				log.Printf("%v", err)
			}
			return
		}
	}

}
