package app

import (
	"fsanalyze/pkg/analyze"

	"github.com/integrii/flaggy"
)

const Version = "1.0.0"

var (
	parse *flaggy.Subcommand
)

var (
	filePath string
)

func init() {
	flaggy.SetName("fsa")
	flaggy.SetDescription("fsanalyze binary parses the fsimage and prints the analyzed info about the hadoop filesystem")

	parse = flaggy.NewSubcommand("parse")
	parse.Description = "parses the given fsimage for analyzed output"
	parse.String(&filePath, "f", "filepath", "filepath for the fsimage")
	flaggy.AttachSubcommand(parse, 1)

	flaggy.SetVersion(Version)
}

func Run() error {
	flaggy.Parse()
	switch {
	// Parse the given fsimage
	case parse.Used:
		if err := analyze.Parse(filePath); err != nil {
			return err
		}
	}
	return nil
}
