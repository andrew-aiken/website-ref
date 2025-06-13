package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/cmd"
	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/flags"
	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/logger"
)

var (
	ProjectName = "cloudfront-frontend"
	Version     = "0.0.1"
)

func main() {
	app := &cli.Command{
		Name:    ProjectName,
		Version: Version,
		Usage:   "Build Vite application and deploy SPA to AWS CloudFront",

		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// Initialize logger with debug mode from flag
			logger.Init(logger.Config{
				Debug: cmd.Bool("debug"),
			})
			return ctx, nil
		},

		Commands: []*cli.Command{
			{
				Name:   "build",
				Usage:  "Installs and builds a Vite application",
				Action: cmd.Build,
				Flags: []cli.Flag{
					flags.OutputDirectory,
					flags.ProjectName,
					flags.ReleaseStage,
					flags.WorkspaceDirectory,
				},
			},
			{
				Name:   "invalidate",
				Usage:  "Invalidates the path to the cloudfront frontend application",
				Action: cmd.InvalidateCloudFront,
				Flags: []cli.Flag{
					flags.ProjectName,
					flags.ReleaseStage,
				},
			},
			{
				Name:   "parse",
				Usage:  "Reads the configuration file and converts it into a default.json for the frontend to use",
				Action: cmd.ParseConfigWrapper,
				Flags: []cli.Flag{
					flags.ConfigPath,
					flags.EnvironmentConfigPath,
					flags.ProjectName,
					flags.ReleaseStage,
					flags.OutputDirectory,
				},
			},
			{
				Name:   "publish",
				Usage:  "Builds, configures, uploads, and invalidates the cloudfront frontend application",
				Action: cmd.Publish,
				Flags: []cli.Flag{
					flags.ConfigPath,
					flags.EnvironmentConfigPath,
					flags.OutputDirectory,
					flags.ProjectName,
					flags.ReleaseStage,
					flags.WorkspaceDirectory,
				},
			},
			{
				Name:   "upload",
				Usage:  "Replaces the existing files in the S3 bucket with the build output",
				Action: cmd.Upload,
				Flags: []cli.Flag{
					flags.EnvironmentConfigPath,
					flags.OutputDirectory,
					flags.ProjectName,
					flags.ReleaseStage,
				},
			},
		},
		Flags: []cli.Flag{
			flags.Debug,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log := logger.Get()
		log.Fatal().Err(err).Msg("Error")
	}
}
