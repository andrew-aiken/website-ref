package cmd

import (
	"context"

	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/logger"
	"github.com/urfave/cli/v3"
)

func Publish(ctx context.Context, cli *cli.Command) error {
	log := logger.Get()

	if err := Build(ctx, cli); err != nil {
		log.Error().Err(err).Msg("Failed to build the project")
		return err
	}

	ctx, err := ParseConfig(ctx, cli)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate the configuration file")
		return err
	}

	if err := Upload(ctx, cli); err != nil {
		log.Error().Err(err).Msg("Failed to upload files to S3")
		return err
	}

	if err := InvalidateCloudFront(ctx, cli); err != nil {
		log.Error().Err(err).Msg("Failed to invalidated CloudFront cache")
		return err
	}

	return nil
}
