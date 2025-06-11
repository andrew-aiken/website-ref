package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/logger"
)

func Upload(ctx context.Context, cli *cli.Command) error {
	// Get source directory and prefix arguments
	project := cli.String("project")
	releaseStage := cli.String("release-stage")
	sourceDir := cli.String("output")

	bucketName, ok := ctx.Value("bucket").(string)
	if !ok {
		// Try parsing the config
		if envMapping, err := readEnvironmentMapping(releaseStage, cli.String("environment-config")); err != nil {
			return fmt.Errorf("bucket not found in context")
		} else {
			bucketName = envMapping.Bucket
		}
	}

	log := logger.Get()

	log.Debug().Msgf("Using bucket: %s", bucketName)

	prefix := filepath.Join(releaseStage, project)

	log.Debug().Msgf("Using prefix: %s", prefix)

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load AWS configuration")
		return err
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	if err := tagObjects(s3Client, ctx, bucketName, prefix); err != nil {
		log.Error().Msgf("Failed to delete existing objects")
		return err
	}

	if err := uploadObjects(s3Client, ctx, bucketName, prefix, sourceDir); err != nil {
		log.Error().Msgf("Failed to upload objects")
		return err
	}

	return nil
}

// tagObjects tags all existing objects in the project prefix with "expire=true".
// This is useful for marking objects for expiration in S3 lifecycle policies.
func tagObjects(s3Client *s3.Client, ctx context.Context, bucketName string, prefix string) error {
	// List all objects with the prefix
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix + "/"), // Add trailing slash to ensure exact prefix match
	}

	paginator := s3.NewListObjectsV2Paginator(s3Client, listInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			_, err := s3Client.PutObjectTagging(ctx, &s3.PutObjectTaggingInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
				Tagging: &types.Tagging{
					TagSet: []types.Tag{
						{
							Key:   aws.String("expire"),
							Value: aws.String("true"),
						},
					},
				},
			})

			if err != nil {
				return fmt.Errorf("failed to tag object %s: %w", *obj.Key, err)
			}
			log.Info().Str("key", *obj.Key).Msg("Tagged object")
		}
	}

	log.Info().Msg("Successfully tagged all objects with prefix")
	return nil
}

func uploadObjects(s3Client *s3.Client, ctx context.Context, bucketName string, prefix string, sourceDir string) error {
	// Walk through the source directory
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Create S3 key with temporary prefix
		s3Key := filepath.Join(prefix, relPath)

		// Determine content type based on file extension
		contentType := "application/octet-stream" // default content type
		switch filepath.Ext(path) {
		case ".html":
			contentType = "text/html"
		case ".css":
			contentType = "text/css"
		case ".js":
			contentType = "application/javascript"
		case ".json":
			contentType = "application/json"
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".gif":
			contentType = "image/gif"
		case ".svg":
			contentType = "image/svg+xml"
		case ".ico":
			contentType = "image/x-icon"
		case ".woff":
			contentType = "font/woff"
		case ".woff2":
			contentType = "font/woff2"
		case ".ttf":
			contentType = "font/ttf"
		case ".eot":
			contentType = "application/vnd.ms-fontobject"
		case ".otf":
			contentType = "font/otf"
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		// Upload file to S3
		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(s3Key),
			Body:        file,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			return fmt.Errorf("failed to upload file %s: %w", path, err)
		}

		log.Info().Str("file", path).Str("s3Key", s3Key).Msg("Uploaded file to S3")
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to upload files to S3")
		return err
	}

	log.Info().Msg("Successfully uploaded all files to S3")
	return nil
}
