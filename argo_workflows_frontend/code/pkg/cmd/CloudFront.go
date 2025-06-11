package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/urfave/cli/v3"

	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/logger"
)

func InvalidateCloudFront(ctx context.Context, cli *cli.Command) error {
	// Get project name and release stage from cli arguments
	project := cli.String("project")
	releaseStage := cli.String("release-stage")

	distributionName, ok := ctx.Value("distribution").(string)
	if !ok {
		return fmt.Errorf("distribution not found in context")
	}

	log := logger.Get()

	log.Debug().Msgf("Using distribution: %s", distributionName)

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load AWS configuration")
		return err
	}

	// Create CloudFront client
	cloudfrontClient := cloudfront.NewFromConfig(cfg)

	// Create invalidation paths
	paths := []string{
		fmt.Sprintf("/%s/index.html", project),   // Invalidate index.html
		fmt.Sprintf("/%s/default.json", project), // Invalidate default.json
	}

	// Create invalidation input
	input := &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distributionName),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(fmt.Sprintf("%s-%s-%d", project, releaseStage, time.Now().Unix())),
			Paths: &types.Paths{
				Quantity: aws.Int32(int32(len(paths))),
				Items:    paths,
			},
		},
	}

	// Create the invalidation
	result, err := cloudfrontClient.CreateInvalidation(ctx, input)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create CloudFront invalidation")
		return err
	}

	log.Info().
		Str("invalidationId", *result.Invalidation.Id).
		Str("status", string(*result.Invalidation.Status)).
		Msg("Successfully created CloudFront invalidation")

	return nil
}
