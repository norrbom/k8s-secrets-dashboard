locals {
  project_name = "k8s-secrets-dashboard.mycompany.com"
}

# inputs are inherited in submodules via include find_in_parent_folders()
inputs = {
  project_name  = local.project_name
}

remote_state {
  backend = "s3"
  config = {
    bucket         = "terraform-state-${local.project_name}"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = "eu-north-1"
    encrypt        = true
    dynamodb_table = "terraform-lock-table-${local.project_name}"
  }
}

generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite"
  contents = <<EOF
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
  required_version = ">= 0.14.9"

  backend "s3" {}
}

provider "aws" {
  profile = "default"
  region  = "eu-north-1"
}
EOF
}