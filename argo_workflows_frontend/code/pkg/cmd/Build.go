package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/andrew-aiken/website-ref/argo_workflows_frontend/code/pkg/logger"
	"github.com/urfave/cli/v3"
)

func Build(ctx context.Context, cli *cli.Command) error {
	log := logger.Get()

	if err := installPackages(cli.String("workspace")); err != nil {
		log.Error().Err(err).Msg("Failed to install packages")
		return err
	}

	if err := buildProject(cli.String("workspace"), cli.String("output")); err != nil {
		log.Error().Err(err).Msg("Failed to install packages")
		return err
	}

	return nil
}

func installPackages(workspace string) error {
	log := logger.Get()

	log.Debug().Msg("Pulling Node modules")

	commandInstall := exec.Command("/usr/local/bin/npm", "clean-install", "--omit=dev", "--ignore-scripts")
	commandInstall.Dir = workspace

	commandInstall.Stdout = os.Stdout
	commandInstall.Stderr = os.Stderr

	if err := commandInstall.Run(); err != nil {
		return err
	}
	return nil
}

func buildProject(workspace string, outputDir string) error {
	log := logger.Get()

	log.Debug().Msg("Building Vite application")

	commandBuild := exec.Command("/usr/local/bin/node", "node_modules/vite/bin/vite.js", "build", "--outDir", outputDir, "--emptyOutDir")
	commandBuild.Dir = workspace

	commandBuild.Stdout = os.Stdout
	commandBuild.Stderr = os.Stderr

	if err := commandBuild.Run(); err != nil {
		return err
	}
	return nil
}
