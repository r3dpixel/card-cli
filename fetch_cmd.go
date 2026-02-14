package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/r3dpixel/card-fetcher/task"
	"github.com/r3dpixel/card-parser/png"
	"github.com/r3dpixel/toolkit/filex"
	"github.com/r3dpixel/toolkit/stringsx"
	"github.com/r3dpixel/toolkit/templater"
	"github.com/schollz/progressbar/v3"
)

var fileNameTemplater = templater.New(tokens...)

func handleFetch(urls []string, output string, format string, chromePath string) error {
	var err error
	if stringsx.IsBlank(output) {
		output, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	if !filex.DirExists(output) {
		return fmt.Errorf("output directory '%s' doesn't exist", output)
	}

	if stringsx.IsBlank(format) {
		format = `{{SOURCE}}_{{PLATFORM_ID}}`
	}
	template := fileNameTemplater.Compile(format)

	r := initializeRouter(chromePath)
	tasks := r.TaskSliceOf(urls...)

	fmt.Println(ImportantTextStyle.Render(fmt.Sprintf("Fetching %d URLs", len(urls))))
	bar := progressbar.Default(int64(len(urls)), "Fetching cards...")
	bar.Set(len(tasks.InvalidURLs))

	failedURLs := make([]string, 0)
	successURLs := make([]string, 0)
	for _, fetchTask := range tasks.Tasks {
		success := executeFetchTask(fetchTask, output, template)
		if success {
			successURLs = append(successURLs, fetchTask.OriginalURL())
		} else {
			failedURLs = append(failedURLs, fetchTask.OriginalURL())
		}
		bar.Add(1)
	}

	printReport(urls, successURLs, failedURLs, tasks.InvalidURLs)

	listURLs("Success URLs:", successURLs, true, SuccessTextStyle)
	listURLs("Failed URLs:", failedURLs, true, DangerTextStyle)
	listURLs("Invalid URLs:", tasks.InvalidURLs, false, ImportantTextStyle)

	return nil
}

func executeFetchTask(t task.Task, output string, template *compiledTemplate) bool {
	metadata, card, err := t.FetchAll()
	if err != nil || card.Sheet == nil || !card.Integrity() || !metadata.IsConsistentWith(card.Sheet) {
		return false
	}

	rawCard, err := card.Encode()
	if err != nil {
		return false
	}

	err = rawCard.ToFile(filepath.Join(output, filex.SanitizePath(template.Execute(metadata))+png.Extension))
	if err != nil {
		return false
	}

	return true
}

func printReport(urls, successURLs, failedURLs, invalidURLs []string) {
	noTotalURLs := InfoTextStyle.Render(fmt.Sprintf("%d Total", len(urls)))
	noSuccessURLs := SuccessTextStyle.Render(fmt.Sprintf("%d Success", len(successURLs)))
	noFailedURLs := DangerTextStyle.Render(fmt.Sprintf("%d Failed", len(failedURLs)))
	noInvalidURLs := ImportantTextStyle.Render(fmt.Sprintf("%d Invalid", len(invalidURLs)))

	report := lipgloss.JoinHorizontal(lipgloss.Top,
		noTotalURLs,
		SubtleTextStyle.Render(" | "),
		noSuccessURLs,
		SubtleTextStyle.Render(" | "),
		noFailedURLs,
		SubtleTextStyle.Render(" | "),
		noInvalidURLs,
	)
	fmt.Println(report)
}

func listURLs(label string, urls []string, endLine bool, style lipgloss.Style) {
	if len(urls) == 0 {
		return
	}

	fmt.Println(style.Render(label))
	for _, url := range urls {
		fmt.Println(style.Render(url))
	}
	if endLine {
		fmt.Println()
	}
}
