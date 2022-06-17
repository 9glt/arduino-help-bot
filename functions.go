package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/bwmarrin/discordgo"
)

func fnHelp(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	_, err := ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Hello World! %s", s))
	if err != nil {
		log.Printf("%v", err)
	}
}

func fnReload(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {

	if !checkUserRole(dm.Member.Roles) {
		log.Printf("User %s is not allowed to use this command", dm.Author.Username)
		return
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = envDocsDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("git pull error: %v ( %s )", err, output)
	}

	go ScanForTags()
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
