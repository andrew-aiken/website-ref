package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func main() {
	awsEcrEndpoint := os.Getenv("AWS_ECR_ENDPOINT")
	tarFile := os.Getenv("IMAGE_TAR")
	EcrImageUri := os.Getenv("IMAGE_URI")

	if awsEcrEndpoint == "" || tarFile == "" || EcrImageUri == "" {
		log.Fatalf("environment variables AWS_ECR_ENDPOINT, IMAGE_TAR, and IMAGE_URI must be set")
	}

	// Load the default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS configuration: %v", err)
	}

	// Create an ECR client
	ecrClient := ecr.NewFromConfig(cfg)

	// Call GetAuthorizationToken to get the ECR login token
	resp, err := ecrClient.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		log.Fatalf("failed to get ECR authorization token: %v", err)
	}

	// Decode the token and extract the Docker login credentials
	for _, authData := range resp.AuthorizationData {
		decodedToken, err := base64.StdEncoding.DecodeString(*authData.AuthorizationToken)
		if err != nil {
			log.Fatalf("failed to decode authorization token: %v", err)
		}

		// The token is in the format "username:password"
		credentials := strings.SplitN(string(decodedToken), ":", 2)
		if len(credentials) != 2 {
			log.Fatalf("unexpected token format: %s", string(decodedToken))
		}

		username := credentials[0]
		password := credentials[1]

		fmt.Println("Loading tarball image...")
		img, err := tarball.ImageFromPath(tarFile, nil)
		if err != nil {
			log.Fatalf("Failed to load tarball image: %v", err)
		}

		// Set up authentication with username and password.
		auth := &authn.Basic{
			Username: username,
			Password: password,
		}

		// Push the image to the repository.
		fmt.Println("Pushing image to repository...")
		err = crane.Push(img, EcrImageUri, crane.WithAuth(auth))
		if err != nil {
			log.Fatalf("Failed to push image to repository: %v", err)
		}

		fmt.Println("Image successfully pushed to", EcrImageUri)
	}
}
