package flags

import "github.com/urfave/cli/v3"

var (
	// Debug defines the debug flag
	Debug = &cli.BoolFlag{
		Name:     "debug",
		Aliases:  []string{"d"},
		OnlyOnce: true,
		Required: false,
		Usage:    "Display debug information",
	}

	ProjectName = &cli.StringFlag{
		Name:     "project",
		OnlyOnce: true,
		Required: true,
		Usage:    "Name of the project",
	}

	ReleaseStage = &cli.StringFlag{
		Name:     "release-stage",
		Aliases:  []string{"stage"},
		OnlyOnce: true,
		Required: true,
		Usage:    "The release stage to deploy to",
	}
)
