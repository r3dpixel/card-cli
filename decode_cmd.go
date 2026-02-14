package main

import (
	"fmt"

	"github.com/r3dpixel/card-parser/png"
	"github.com/r3dpixel/toolkit/jsonx"
	"github.com/r3dpixel/toolkit/sonicx"
	"github.com/r3dpixel/toolkit/stringsx"
)

type decodeOptions struct {
	pretty     bool
	stable     bool
	outputFile string
}

func handleDecode(inputFile string, opts decodeOptions) error {
	if opts.stable {
		sonicx.Config = sonicx.StableSort
	}

	rawCard, err := png.FromFile(inputFile).LastVersion().Get()
	if err != nil {
		return err
	}
	editableCard, err := rawCard.Decode()
	if err != nil {
		return err
	}

	if stringsx.IsBlank(opts.outputFile) {
		jsonData, err := editableCard.ToBytes()
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))
		return nil
	}

	return editableCard.Sheet.ToFile(opts.outputFile, jsonx.Options{
		Pretty: opts.pretty,
		Indent: "    ",
	})
}
