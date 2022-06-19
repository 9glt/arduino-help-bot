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
	envGuildID       = os.Getenv("BOT_GUILD_ID")

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

	var tags []*discordgo.ApplicationCommandOption

	for tag, v := range tagsRegistry {
		_tag := &discordgo.ApplicationCommandOption{
			Name:        strings.ToLower(tag),
			Description: v.Title,
			Type:        discordgo.ApplicationCommandOptionSubCommand,
		}
		tags = append(tags, _tag)
		log.Printf("command: %v :%v", tag, v.Title)
	}

	commands = append(commands, &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "help tags",
		// Options:     tags,
	})

	commands = append(commands, &discordgo.ApplicationCommand{
		Name:        "tag",
		Description: "help tags",
		Options:     tags,
	})

	log.Println("Adding commands...")

	commandsInServer, _ := dg.ApplicationCommands(dg.State.User.ID, envGuildID)
	commandsInServerMap := make(map[string]struct{})
	for _, cmd := range commandsInServer {
		commandsInServerMap[cmd.Name] = struct{}{}
		log.Printf("%v", cmd.Name)
		// dg.ApplicationCommandDelete(dg.State.User.ID, "275273435930951680", cmd.ID)
	}
	// return
	// dg.ApplicationCommands(dg.State.User.ID)
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	// commands = commands[len(commands)-2:]
	for i, v := range commands {
		if _, ok := commandsInServerMap[v.Name]; ok {
			// _, err := dg.ApplicationCommandEdit(dg.State.User.ID, "275273435930951680", v.ID, v)
			// if err != nil {
			// 	log.Printf("Cannot create '%v' command: %v", v.Name, err)
			// }
			log.Printf("Command already in server: %v", v.Name)
			continue
		}
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, envGuildID, v)
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
