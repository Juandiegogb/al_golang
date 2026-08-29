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
    archive = {
      source  = "hashicorp/archive"
      version = "2.8.0"
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

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "example" {
  name               = "lambda_execution_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "${path.module}/bootstrap"
  output_path = "/tmp/function.zip"
}

resource "aws_lambda_function" "utils_function" {
  function_name    = "utils_lambda_function"
  architectures    = ["arm64"]
  role             = aws_iam_role.example.arn
  filename         = data.archive_file.zip.output_path
  runtime          = "provided.al2023"
  handler          = "main.handler"
  source_code_hash = data.archive_file.zip.output_base64sha256
}

resource "aws_lambda_function_url" "function_url" {
  authorization_type = "NONE"
  function_name      = aws_lambda_function.utils_function.function_name
}

output "function_url" {
  value = aws_lambda_function_url.function_url.function_url
}
