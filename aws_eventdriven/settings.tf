terraform {
  required_version = "1.9.6"

  backend "local" {
    path = "./terraform.tfstate"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.69.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = {
      terraform = "true"
      purpose   = "c3_eventdriven"
    }
  }
}
