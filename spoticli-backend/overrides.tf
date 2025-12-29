# Local override to use LocalStack for ALL AWS services
# This lets you practice Terraform without touching real AWS
provider "aws" {
  region                      = var.aws_region
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  s3_use_path_style           = true

  endpoints {
    s3             = "http://localhost:4566"
    rds            = "http://localhost:4566"
    ecs            = "http://localhost:4566"
    ec2            = "http://localhost:4566"
    iam            = "http://localhost:4566"
    cloudwatch     = "http://localhost:4566"
    cloudwatchlogs = "http://localhost:4566"
  }
}

# Note: To use REAL AWS, simply rename this file (e.g., overrides.tf.disabled)
# and Terraform will use the default AWS provider from main.tf
