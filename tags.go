package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type TagField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool
}

type Tag struct {
	Title   string     `json:"title"`
	Aliases []string   `json:"aliases"`
	Fields  []TagField `json:"fields"`
	Image   string     `json:"image"`
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
	for _, alias := range tag.Aliases {
		tagsRegistry[alias] = &tag
	}

	return nil
}

func ScanForTags() {
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
