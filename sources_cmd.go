package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/r3dpixel/card-fetcher/router"
	"github.com/r3dpixel/card-fetcher/source"
)

type SourceStyle struct {
	BgColor string
	FgColor string
}

var styles = map[source.ID]SourceStyle{
	source.CharacterTavern: {BgColor: "#e83e8c", FgColor: "#FFFFFF"},
	source.ChubAI:          {BgColor: "#e5e7eb", FgColor: "#1f2937"},
	source.NyaiMe:          {BgColor: "#6b7280", FgColor: "#FFFFFF"},
	source.PepHop:          {BgColor: "#0d9488", FgColor: "#FFFFFF"},
	source.Pygmalion:       {BgColor: "#4b0082", FgColor: "#FFFFFF"},
	source.WyvernChat:      {BgColor: "#f97316", FgColor: "#FFFFFF"},
	source.JannyAI:         {BgColor: "#55a2fa", FgColor: "#FFFFFF"},
	source.AICC:            {BgColor: "#753fba", FgColor: "#FFFFFF"},
}

func listSources(pretty bool) {
	r := initializeRouter("")
	sources := r.Sources()
	statuses := make(map[source.ID]router.IntegrationStatus)
	for _, s := range sources {
		statuses[s] = r.CheckIntegration(s)
	}

	if !pretty {
		listSourcesSimple(r, statuses)
		return
	}

	maxContentWidth := 0
	for _, f := range r.Fetchers() {
		labelWidth := lipgloss.Width(string(f.SourceID()))
		urlWidth := lipgloss.Width(f.SourceURL())
		statusWidth := lipgloss.Width(string(statuses[f.SourceID()]))
		maxContentWidth = max(maxContentWidth, urlWidth, statusWidth, labelWidth)
	}

	var styledSources []string
	for _, f := range r.Fetchers() {
		sourceLabel := string(f.SourceID())
		sourceURL := f.SourceURL()
		status := statuses[f.SourceID()]

		finalLineWidth := maxContentWidth + 2

		mainStyle := lipgloss.NewStyle().
			Background(lipgloss.Color(styles[f.SourceID()].BgColor)).
			Foreground(lipgloss.Color(styles[f.SourceID()].FgColor))

		var statusStyle lipgloss.Style
		if status == router.IntegrationSuccess {
			statusStyle = lipgloss.NewStyle().
				Foreground(ColorNeutral).
				Background(ColorSuccess).
				Bold(true)
		} else {
			statusStyle = lipgloss.NewStyle().
				Foreground(ColorNeutral).
				Background(ColorDanger).
				Bold(true)
		}

		styledLabel := mainStyle.Render(
			lipgloss.PlaceHorizontal(finalLineWidth, lipgloss.Center, sourceLabel),
		)
		styledURL := mainStyle.Render(
			lipgloss.PlaceHorizontal(finalLineWidth, lipgloss.Center, sourceURL),
		)

		styledStatus := statusStyle.Render(
			lipgloss.PlaceHorizontal(finalLineWidth, lipgloss.Center, string(status)),
		)

		finalTagContent := lipgloss.JoinVertical(lipgloss.Top,
			styledLabel,
			styledURL,
			styledStatus,
		)

		containerStyle := lipgloss.NewStyle().MarginRight(2)
		renderedTag := containerStyle.Render(finalTagContent)
		styledSources = append(styledSources, renderedTag)
	}

	var rows []string
	rowStyle := lipgloss.NewStyle().MarginBottom(1)
	for i := 0; i < len(styledSources); i += 3 {
		end := min(i+3, len(styledSources))
		row := lipgloss.JoinHorizontal(lipgloss.Top, styledSources[i:end]...)
		if end >= len(styledSources) {
			rows = append(rows, row)
		} else {
			rows = append(rows, rowStyle.Render(row))
		}
	}

	finalLayout := lipgloss.JoinVertical(lipgloss.Left, rows...)
	fmt.Println(finalLayout)
}

func listSourcesSimple(r *router.Router, statuses map[source.ID]router.IntegrationStatus) {
	for _, f := range r.Fetchers() {
		sourceLabel := string(f.SourceID())
		status := statuses[f.SourceID()]

		sourceStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles[f.SourceID()].BgColor)).
			Bold(true)

		var statusStyle lipgloss.Style
		if status == router.IntegrationSuccess {
			statusStyle = SuccessTextStyle
		} else {
			statusStyle = DangerTextStyle
		}

		fmt.Printf("%s: %s\n", sourceStyle.Render(sourceLabel), statusStyle.Render(string(status)))
	}
}
