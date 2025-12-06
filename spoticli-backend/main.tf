terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

# AWS Provider Configuration
provider "aws" {
  region                   = var.aws_region
  shared_config_files      = ["~/.aws/config"]
  shared_credentials_files = ["~/.aws/credentials"]
  profile                  = "matttm"
}

# Variables
variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-east-1"
}

variable "app_port" {
  description = "Application port"
  type        = number
  default     = 4200
}

variable "fargate_cpu" {
  description = "Fargate instance CPU units (256, 512, 1024, 2048, 4096)"
  type        = number
  default     = 256
}

variable "fargate_memory" {
  description = "Fargate instance memory in MB (512, 1024, 2048, etc.)"
  type        = number
  default     = 512
}

variable "app_image" {
  description = "Full image reference (e.g. 123456789012.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest or dockerhub-user/spoticli-backend:latest)"
  type        = string
}

# Data source for current AWS account
data "aws_caller_identity" "current" {}

# Use default VPC (already exists in your AWS account)
data "aws_vpc" "default" {
  default = true
}

# Get default subnets (already exist)
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# S3 Bucket for Music Storage (simple configuration)
resource "aws_s3_bucket" "music_storage" {
  bucket = "spoticli-music-${data.aws_caller_identity.current.account_id}"

  tags = {
    Name = "Spoticli Music Storage"
  }
}

# Security Group for Fargate Tasks
resource "aws_security_group" "fargate_tasks" {
  name        = "spoticli-fargate-sg"
  description = "Allow access to Spoticli app"
  vpc_id      = data.aws_vpc.default.id

  # Application port
  ingress {
    description = "Application Port"
    from_port   = var.app_port
    to_port     = var.app_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "spoticli-fargate-sg"
  }
}

# IAM Role for ECS Task Execution (required by Fargate)
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "spoticli-ecs-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# IAM Role for the application (to access S3)
resource "aws_iam_role" "ecs_task_role" {
  name = "spoticli-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ecs-tasks.amazonaws.com"
      }
    }]
  })
}

# Allow app to access S3 bucket
resource "aws_iam_role_policy" "s3_access" {
  name = "s3-access"
  role = aws_iam_role.ecs_task_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ]
      Resource = [
        aws_s3_bucket.music_storage.arn,
        "${aws_s3_bucket.music_storage.arn}/*"
      ]
    }]
  })
}

# CloudWatch Log Group (to see your app logs)
resource "aws_cloudwatch_log_group" "spoticli_backend" {
  name              = "/ecs/spoticli-backend"
  retention_in_days = 3

  tags = {
    Name = "spoticli-logs"
  }
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "spoticli-cluster"

  tags = {
    Name = "spoticli-cluster"
  }
}

# ECS Task Definition (describes your container)
resource "aws_ecs_task_definition" "spoticli_backend" {
  family                   = "spoticli-backend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.fargate_cpu
  memory                   = var.fargate_memory
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([{
    name      = "spoticli-backend"
    image     = var.app_image
    essential = true

    portMappings = [{
      containerPort = var.app_port
      protocol      = "tcp"
    }]

    environment = [
      {
        name  = "MUSIC_BUCKET"
        value = aws_s3_bucket.music_storage.bucket
      },
      {
        name  = "PORT"
        value = tostring(var.app_port)
      }
    ]

    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.spoticli_backend.name
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])
}

# ECS Service (runs your container)
resource "aws_ecs_service" "spoticli_backend" {
  name            = "spoticli-backend-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.spoticli_backend.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.fargate_tasks.id]
    subnets          = data.aws_subnets.default.ids
    assign_public_ip = true
  }

  tags = {
    Name = "spoticli-service"
  }
}

# Outputs - Information you'll need after deployment
output "music_bucket_name" {
  description = "S3 bucket for music files"
  value       = aws_s3_bucket.music_storage.bucket
}

output "get_app_url" {
  description = "Run this command to get your app URL"
  value       = "./get-task-ip.sh"
}
