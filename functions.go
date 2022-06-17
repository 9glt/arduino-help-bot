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

func fnReload(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	log.Printf("%v %v", dm.Member.Roles, checkUserRole(dm.Member.Roles))
	ScanForTags()
}

func fnTag(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	tagsRegistryLock.RLock()
	tag, ok := tagsRegistry[s]
	tagsRegistryLock.RUnlock()

	if !ok {
		return
	}

	_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, &discordgo.MessageEmbed{
		Title: tag.Title,
		Image: &discordgo.MessageEmbedImage{
			URL: tag.Image,
		},
		Fields: tag.Fields,
	})

	if err != nil {
		log.Printf("%v", err)
	}
}
