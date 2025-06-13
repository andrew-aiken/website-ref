package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"

	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/logger"
	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/types"
)

// ParseConfigWrapper is a wrapper function that matches the cli.ActionFunc signature
func ParseConfigWrapper(ctx context.Context, cli *cli.Command) error {
	_, err := ParseConfig(ctx, cli)
	return err
}

func ParseConfig(ctx context.Context, cli *cli.Command) (context.Context, error) {
	configPath := cli.String("config")
	envConfigPath := cli.String("environment-config")
	output := cli.String("output")
	project := cli.String("project")
	releaseStage := cli.String("release-stage")

	log := logger.Get()

	projectConfig, err := readConfigFile(releaseStage, project, configPath)
	if err != nil {
		return ctx, err
	}

	log.Info().Msgf("Parsed config: %+v", projectConfig)

	envMapping, err := readEnvironmentMapping(releaseStage, envConfigPath)
	if err != nil {
		return ctx, err
	}

	// Set the bucket & distribution in the context to be used later
	ctx = context.WithValue(ctx, "bucket", envMapping.Bucket)
	ctx = context.WithValue(ctx, "distribution", envMapping.Distribution)

	bucket, ok := ctx.Value("bucket").(string)
	if !ok {
		return ctx, fmt.Errorf("bucket not found in context")
	}
	log.Info().Msgf("Using bucket: %s", bucket)

	log.Info().Msgf("Environment: %s, Domain: %s", releaseStage, envMapping.Domain)

	if err := writeConfig(releaseStage, envMapping.Domain, projectConfig, output); err != nil {
		log.Error().Err(err).Msg("Failed to write config")
		return ctx, err
	}

	return ctx, nil
}

func readConfigFile(releaseStage string, project string, filePath string) (map[string]string, error) {
	// Read the configuration file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
		return nil, err
	}

	var cfg types.Stages
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal yaml")
		return nil, err
	}

	log.Debug().Msgf("Parsed config: %+v", cfg)

	projectConfig, ok := cfg[releaseStage][project]
	if !ok {
		return nil, fmt.Errorf("No config found for project %s in stage %s", project, releaseStage)
	}
	return projectConfig, nil
}

// readEnvironmentMapping reads the environment mapping from the specified file path and returns the configuration for the given release stage
func readEnvironmentMapping(releaseStage string, filePath string) (types.EnvironmentConfig, error) {
	// Read the environment yaml file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read config file")
		return types.EnvironmentConfig{}, err
	}

	// Unmarshal the yaml data into the environment mapping structure
	var envMapping types.EnvironmentMapping
	err = yaml.Unmarshal(data, &envMapping)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal yaml")
		return types.EnvironmentConfig{}, err
	}

	log.Debug().Msgf("Parsed env mapping: %+v", envMapping)

	// Check if the release stage exists in the mapping
	environments, ok := envMapping[releaseStage]
	if !ok {
		return types.EnvironmentConfig{}, fmt.Errorf("Invalid release stage: %s", releaseStage)
	}

	return environments, nil
}

func writeConfig(releaseStage string, domain string, projectConfig map[string]string, outputDir string) error {
	log.Debug().Msgf("Writing config for release stage: %s, domain: %s", releaseStage, domain)

	configTemplate := types.ProjectConfig{
		ReleaseStage: releaseStage,
		Urls: map[string]string{
			"apiUrl":   fmt.Sprintf("https://api.%s", domain),
			"appUrl":   fmt.Sprintf("https://app.%s", domain),
			"loginUrl": fmt.Sprintf("https://login.%s", domain),
		},
	}
	mergedMap, err := mergeToMap(configTemplate, projectConfig)
	if err != nil {
		panic(err)
	}

	// Marshal the merged map to JSON and add indentation
	jsonBytes, err := json.MarshalIndent(mergedMap, "", "  ")
	if err != nil {
		panic(err)
	}

	// Write JSON to a file
	outputFile := filepath.Join(outputDir, "default.json")
	err = os.WriteFile(outputFile, jsonBytes, 0644)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to write JSON to file: %s", outputFile)
		return err
	}

	log.Debug().Msg("JSON written to file")
	log.Debug().Msg(string(jsonBytes))

	return nil
}

func mergeToMap(a, b interface{}) (map[string]interface{}, error) {
	mapA := make(map[string]interface{})
	mapB := make(map[string]interface{})

	// Marshal and unmarshal to maps
	aBytes, _ := yaml.Marshal(a)
	bBytes, _ := yaml.Marshal(b)

	yaml.Unmarshal(aBytes, &mapA)
	yaml.Unmarshal(bBytes, &mapB)

	for k, v := range mapB {
		mapA[k] = v
	}

	return mapA, nil
}
