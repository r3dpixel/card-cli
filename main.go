package main

import (
	"errors"
	"os"
	"strings"

	"github.com/r3dpixel/card-fetcher/impl"
	"github.com/r3dpixel/card-fetcher/source"
	"github.com/r3dpixel/toolkit/trace"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func init() {
	zerolog.ErrorMarshalFunc = trace.ErrorMarshalFunc
	log.Logger = log.Logger.Output(trace.ConsoleTraceWriter()).Level(zerolog.Disabled)
}

var sources = []source.ID{
	source.CharacterTavern,
	source.ChubAI,
	source.NyaiMe,
	source.PepHop,
	source.WyvernChat,
}

func main() {
	app := &cli.App{
		Name:  "card-cli",
		Usage: "A tool for fetching, decoding, and modifying V2/V3 chara cards",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose logging output",
			},
		},
		Before: func(cCtx *cli.Context) error {
			if cCtx.Bool("verbose") {
				log.Logger = log.Logger.Level(zerolog.TraceLevel)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "fetch",
				Usage:     "Fetches cards from one or more URLs",
				UsageText: "card-cli fetch [--format / -f STR] [--output / -o FOLDER] [--autofetch] [--chromepath STR] <URL1> [URL2]...",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "The path where the cards will be saved. If missing current directory will be used.",
					},
					&cli.StringFlag{
						Name:    "format",
						Aliases: []string{"f"},
						Usage:   "The format of the file names (use any of the following macros): " + strings.Join(tokenKeys(), ", "),
					},
					&cli.BoolFlag{
						Name:  "autofetch",
						Usage: "Automatically fetch related cards when available",
					},
					&cli.StringFlag{
						Name:  "chromepath",
						Usage: "Path to Chrome/Chromium executable for JannyAI source",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() == 0 {
						return errors.New("'fetch' command requires at least one argument: an URL or list of URLs")
					}
					output := cCtx.String("output")
					format := cCtx.String("format")

					urls := cCtx.Args().Slice()

					chromeConfigProvider := func() impl.JannyChromeConfig {
						return impl.JannyChromeConfig{
							Path:             cCtx.String("chromepath"),
							AutoFetchCookies: cCtx.Bool("autofetch"),
						}
					}

					return handleFetch(urls, chromeConfigProvider, output, format)
				},
			},
			{
				Name:      "decode",
				Usage:     "Decodes a target chara card and outputs the JSON",
				UsageText: "card-cli decode [--pretty / -p] [--stable / -s] [--output / -o FILE] <file>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Specify output file name. If missing STD_OUT will be used.",
					},
					&cli.BoolFlag{
						Name:    "pretty",
						Aliases: []string{"p"},
						Usage:   "Prettify JSON output",
					},
					&cli.BoolFlag{
						Name:    "stable",
						Aliases: []string{"s"},
						Usage:   "Use stable sort for keys in JSON output",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 1 {
						return errors.New("'decode' command requires exactly one argument: a card file")
					}
					inputFile := cCtx.Args().First()
					outputFile := cCtx.String("output")
					pretty := cCtx.Bool("pretty")
					stable := cCtx.Bool("stable")
					return handleDecode(inputFile, decodeOptions{
						outputFile: outputFile,
						pretty:     pretty,
						stable:     stable,
					})
				},
			},
			{
				Name:      "inject",
				Usage:     "Replaces the JSON data from a chara card",
				UsageText: "card-cli inject <card> <json>",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 2 {
						return errors.New("'inject' command requires exactly two arguments: a card file and a json file")
					}
					imageFile := cCtx.Args().Get(0)
					jsonFile := cCtx.Args().Get(1)
					return handleInject(imageFile, jsonFile)
				},
			},
			{
				Name:      "sources",
				Usage:     "Lists supported external sources (i.e ChubAI)",
				UsageText: "card-cli sources [--pretty / -p] [--autofetch] [--chromepath STR]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "pretty",
						Aliases: []string{"p"},
						Usage:   "Display sources in a colorful dashboard layout",
					},
					&cli.BoolFlag{
						Name:  "autofetch",
						Usage: "Automatically fetch related cards when available",
					},
					&cli.StringFlag{
						Name:  "chromepath",
						Usage: "Path to Chrome/Chromium executable for JannyAI source",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 0 {
						return errors.New("'sources' command requires no arguments")
					}

					chromeConfigProvider := func() impl.JannyChromeConfig {
						return impl.JannyChromeConfig{
							Path:             cCtx.String("chromepath"),
							AutoFetchCookies: cCtx.Bool("autofetch"),
						}
					}

					listSources(chromeConfigProvider, cCtx.Bool("pretty"))
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msg(err.Error())
	}
}
