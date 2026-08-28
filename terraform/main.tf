terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>6.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "~>3.0"
    }
  }

  backend "s3" {
    bucket       = "mts-tf-states-725740881379-us-west-1-an"
    key          = "pl/learning/al-golang/state-file"
    region       = "us-west-1"
    use_lockfile = true
  }
}

provider "aws" {
  region = "us-east-1"
}
