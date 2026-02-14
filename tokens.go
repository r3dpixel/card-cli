package main

import (
	"strconv"

	"github.com/r3dpixel/card-fetcher/models"
	"github.com/r3dpixel/toolkit/templater"
	"github.com/r3dpixel/toolkit/timestamp"
)

type compiledTemplate = templater.CompiledTemplate[*models.Metadata]

type token = templater.Token[*models.Metadata]
type basicToken = templater.BasicToken[*models.Metadata]
type richToken = templater.RichToken[*models.Metadata]

var tokens = []token{
	&richToken{
		BasicToken: basicToken{
			Key: "{{SOURCE}}",
			Extractor: func(metadata *models.Metadata) string {
				return string(metadata.Source)
			},
		},
		Description: "The source platform of the card",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{PLATFORM_ID}}",
			Extractor: func(metadata *models.Metadata) string {
				return metadata.CardInfo.PlatformID
			},
		},
		Description: "The unique platform identifier for the card",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{CHARACTER_ID}}",
			Extractor: func(metadata *models.Metadata) string {
				return metadata.CharacterID
			},
		},
		Description: "The unique character identifier",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{TITLE}}",
			Extractor: func(metadata *models.Metadata) string {
				return metadata.Title
			},
		},
		Description: "The title of the card",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{NAME}}",
			Extractor: func(metadata *models.Metadata) string {
				return metadata.Name
			},
		},
		Description: "The name of the character",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{CREATOR}}",
			Extractor: func(metadata *models.Metadata) string {
				return metadata.Nickname
			},
		},
		Description: "The nickname of the card creator",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{CREATE_TIME}}",
			Extractor: func(metadata *models.Metadata) string {
				return strconv.Itoa(int(timestamp.ConvertToSeconds(metadata.CreateTime)))
			},
		},
		Description: "The creation timestamp in seconds",
	},
	&richToken{
		BasicToken: basicToken{
			Key: "{{UPDATE_TIME}}",
			Extractor: func(metadata *models.Metadata) string {
				return strconv.Itoa(int(timestamp.ConvertToSeconds(metadata.UpdateTime)))
			},
		},
		Description: "The last update timestamp in seconds",
	},
}

func tokenKeys() []string {
	keys := make([]string, 0, len(tokens))
	for _, t := range tokens {
		keys = append(keys, t.GetKey())
	}
	return keys
}
