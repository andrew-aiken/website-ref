# Building and Pushing Container Images to Multiple Registries with Argo Workflows

This is the reference material for my blog post on use Argo Workflows, Kaniko and Crane to push to multiple ECR repositories in different AWS accounts or partitions.

[infrasec.sh/post/argo-workflows-multi-registry](https://infrasec.sh/post/argo-workflows-multi-registry/)

## Usage

## Local Testing

```bash
go build src/main.go

export AWS_ECR_ENDPOINT=0123456789.dkr.ecr.us-east-1.amazonaws.com
export IMAGE_TAR=image.tar
export IMAGE_URI=0123456789.dkr.ecr.us-east-1.amazonaws.com/infrasec-crane:latest

./main
```

## Docker

```bash
docker build --platform=linux/amd64 -t crane .
```