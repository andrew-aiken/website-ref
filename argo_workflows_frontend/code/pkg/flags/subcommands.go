package flags

import "github.com/urfave/cli/v3"

var (
	ConfigPath = &cli.StringFlag{
		Name:     "config",
		OnlyOnce: true,
		Required: false,
		Usage:    "The path to the configuration file",
		Value:    "/mnt/config.yaml",
	}

	EnvironmentConfigPath = &cli.StringFlag{
		Name:     "environment-config",
		OnlyOnce: true,
		Required: false,
		Usage:    "The path to the environment configuration file",
		Value:    "/mnt/envMapping.yaml",
	}

	OutputDirectory = &cli.StringFlag{
		Name:     "output",
		OnlyOnce: true,
		Required: false,
		Usage:    "Path to the outputted build",
		Value:    "/mnt/dist",
	}

	WorkspaceDirectory = &cli.StringFlag{
		Name:     "workspace",
		OnlyOnce: true,
		Required: false,
		Usage:    "Path to the projects source code",
		Value:    "/mnt/workspace",
	}
)
