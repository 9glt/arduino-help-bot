package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Tag struct {
	Title   string                         `json:"title"`
	Aliases []string                       `json:"aliases"`
	Fields  []*discordgo.MessageEmbedField `json:"fields"`
	Image   string                         `json:"image"`
}

func parseAndRegisterTag(path string) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fh.Close()

	var tag Tag
	err = json.NewDecoder(fh).Decode(&tag)
	if err != nil {
		return err
	}
	tagsRegistryLock.Lock()
	for _, alias := range tag.Aliases {
		tagsRegistry[alias] = &tag
	}
	tagsRegistryLock.Unlock()

	return nil
}

func ScanForTags() {
	tagsRegistryLock.Lock()
	tagsRegistry = make(map[string]*Tag)
	tagsRegistryLock.Unlock()

	err := filepath.Walk("/docs", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".json") {
			err := parseAndRegisterTag(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
