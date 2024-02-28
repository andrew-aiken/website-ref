terraform {
  backend "local" {
    path = "./terraform.tfstate"
  }

  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.4.1"
    }
  }
}

resource "local_file" "example" {
  content  = "foo!"
  filename = "./example_file.txt"
}


data "local_file" "source" {
  filename = "source.txt"
}

resource "local_file" "destination" {
  content  = data.local_file.source.content
  filename = "./destination.txt"
}
