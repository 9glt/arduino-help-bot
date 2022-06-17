package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func fnHelp(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	_, err := ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Hello World! %s", s))
	if err != nil {
		log.Printf("%v", err)
	}
}

func fnTag(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	_, err := ds.ChannelMessageSend(dm.ChannelID, "Tag not found")
	if err != nil {
		log.Printf("%v", err)
	}
}
