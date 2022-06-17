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
	tag, ok := tagsRegistry[s]
	if !ok {
		return
	}

	fields := make([]*discordgo.MessageEmbedField, 0)
	for _, field := range tag.Fields {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: field.Inline,
		})
	}

	var embeds []*discordgo.MessageEmbed
	embeds = append(embeds, &discordgo.MessageEmbed{
		Title: tag.Title,
		Image: &discordgo.MessageEmbedImage{
			URL: tag.Image,
		},
		Fields: fields,
	})

	_, err := ds.ChannelMessageSendEmbeds(dm.ChannelID, embeds)
	if err != nil {
		log.Printf("%v", err)
	}
}
