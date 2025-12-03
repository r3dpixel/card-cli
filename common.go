package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/r3dpixel/card-fetcher/fetcher"
	"github.com/r3dpixel/card-fetcher/impl"
	"github.com/r3dpixel/card-fetcher/router"
	"github.com/r3dpixel/toolkit/cred"
	"github.com/r3dpixel/toolkit/reqx"
	"github.com/r3dpixel/toolkit/stringsx"
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

func initializeRouter(chromeConfig func() impl.JannyChromeConfig) *router.Router {
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
	}

	if !skipJannyAI(chromeConfig) {
		builders = append(builders, impl.JannyAIBuilder{
			ChromeConfig: chromeConfig,
			CookieProvider: func() impl.JannyCookies {
				return impl.JannyCookies{
					CloudflareClearance: os.Getenv("JANNY_CF_COOKIE"),
					UserAgent:           os.Getenv("JANNY_USER_AGENT"),
				}
			},
		})
	} else {
		lineOfOrangeText := ImportantTextStyle.Render("JannyAI Source is skipped. Please set JANNY_CF_COOKIE and JANNY_USER_AGENT as environment variables.")
		fmt.Println(lineOfOrangeText)
		fmt.Println()
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

func skipJannyAI(chromeConfig func() impl.JannyChromeConfig) bool {
	// Check if required env vars are set
	hasEnvVars := stringsx.IsNotBlank(os.Getenv("JANNY_CF_COOKIE")) && stringsx.IsNotBlank(os.Getenv("JANNY_USER_AGENT"))
	// Skip if both chrome config is not ok AND no env vars
	return !router.IsChromeOK(chromeConfig) && !hasEnvVars
}
