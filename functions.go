package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func fnHelp(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	tags := &discordgo.MessageEmbed{}

	tags.Fields = append(tags.Fields, &discordgo.MessageEmbedField{
		Name:   "!about",
		Value:  "Get informationa about bot",
		Inline: true,
	})

	tags.Fields = append(tags.Fields, &discordgo.MessageEmbedField{
		Name:   "!help",
		Value:  "Displays help page",
		Inline: true,
	})

	tags.Fields = append(tags.Fields, &discordgo.MessageEmbedField{
		Name:   "!ping",
		Value:  "Get latency of the bot",
		Inline: true,
	})

	tags.Fields = append(tags.Fields, &discordgo.MessageEmbedField{
		Name:   "!tag",
		Value:  "Send a canned response in the channel",
		Inline: true,
	})

	_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, tags)
	if err != nil {
		log.Printf("%v", err)
	}
}

func fnReload(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {

	if !checkUserRole(dm.Member.Roles) {
		log.Printf("User %s is not allowed to use this command", dm.Author.Username)
		return
	}

	go ScanForTags()
}

func fnTag(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	tagsRegistryLock.RLock()
	tag, ok := tagsRegistry[s]
	tagsRegistryLock.RUnlock()

	if !ok {
		printAllTags(ds, dm)
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

func fnFallback(s string, ds *discordgo.Session, dm *discordgo.MessageCreate) {
	tagsRegistryLock.RLock()
	tag, ok := tagsRegistry[s[1:]]
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

func printAllTags(ds *discordgo.Session, dm *discordgo.MessageCreate) {
	tags := &discordgo.MessageEmbed{}

	tagsRegistryLock.RLock()
	for k, v := range tagsRegistry {

		tags.Fields = append(tags.Fields, &discordgo.MessageEmbedField{
			Name:   k,
			Value:  v.Title,
			Inline: true,
		})
	}
	tagsRegistryLock.RUnlock()
	tags.Footer = &discordgo.MessageEmbedFooter{
		Text: "Use !tag <tag> to get info about tag. ( example: !tag pulldown )",
	}
	_, err := ds.ChannelMessageSendEmbed(dm.ChannelID, tags)
	if err != nil {
		log.Printf("%v", err)
	}
}

func fnPing(ds *discordgo.Session, dm *discordgo.MessageCreate) {

	_, err := ds.ChannelMessageSend(dm.ChannelID, fmt.Sprintf("Pong! Rountrip time: %v, API heartbeat: %v", time.Now().Sub(dm.Timestamp), ds.HeartbeatLatency()))
	if err != nil {
		log.Printf("%v", err)
	}
}
