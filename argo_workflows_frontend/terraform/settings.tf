terraform {
  required_version = "1.12.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.99.1"
    }
  }
}


provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = {
      terraform = "true"
    }
  }
}
