package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/r3dpixel/card-fetcher/fetcher"
	"github.com/r3dpixel/card-fetcher/impl"
	"github.com/r3dpixel/card-fetcher/router"
	"github.com/r3dpixel/toolkit/cred"
	"github.com/r3dpixel/toolkit/reqx"
)

var ColorInfo = lipgloss.Color("#1b58ad")
var ColorImportant = lipgloss.Color("#FFA500")
var ColorSuccess = lipgloss.Color("#008a00")
var ColorDanger = lipgloss.Color("#e20017")
var ColorSubtle = lipgloss.Color("#626262")
var ColorNeutral = lipgloss.Color("#ded1cf")

var InfoTextStyle = lipgloss.NewStyle().Foreground(ColorInfo)
var ImportantTextStyle = lipgloss.NewStyle().Foreground(ColorImportant)
var SuccessTextStyle = lipgloss.NewStyle().Foreground(ColorSuccess)
var DangerTextStyle = lipgloss.NewStyle().Foreground(ColorDanger)
var SubtleTextStyle = lipgloss.NewStyle().Foreground(ColorSubtle)

var pygmalionCredReader cred.IdentityReader = cred.NewManager("pygmalion", cred.Env)

func initializeRouter(chromePath string) *router.Router {
	r := router.New(
		reqx.Options{
			RetryCount:        4,
			MinBackoff:        10 * time.Millisecond,
			MaxBackoff:        500 * time.Millisecond,
			DisableKeepAlives: true,
			Impersonation:     reqx.Chrome,
		},
	)

	builders := []fetcher.Builder{
		impl.CharacterTavernBuilder{},
		impl.ChubAIBuilder{},
		impl.NyaiMeBuilder{},
		impl.PephopBuilder{},
		impl.WyvernChatBuilder{},
		impl.AiccBuilder{},
		impl.JannyAIBuilder{ChromePath: func() string { return chromePath }},
	}

	if !skipPygmalion() {
		builders = append(builders, impl.PygmalionBuilder{IdentityReader: pygmalionCredReader})
	} else {
		lineOfOrangeText := ImportantTextStyle.Render("Pygmalion Source is skipped. Please set PYGMALION_USERNAME and PYGMALION_PASSWORD as environment variables.")
		fmt.Println(lineOfOrangeText)
		fmt.Println()
	}

	r.RegisterBuilders(builders...)

	return r
}

func skipPygmalion() bool {
	_, err := pygmalionCredReader.Get()
	return err != nil
}
