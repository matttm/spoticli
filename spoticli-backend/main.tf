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
# Note: Provider is configured in overrides.tf for LocalStack
# For real AWS, remove/rename overrides.tf and uncomment below:
# provider "aws" {
#   region                   = var.aws_region
#   shared_config_files      = ["~/.aws/config"]
#   shared_credentials_files = ["~/.aws/credentials"]
#   profile                  = "matttm"
# }

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
  bucket        = var.tracks_bucket_name
  force_destroy = true

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

# Security Group for the database
resource "aws_security_group" "db" {
  name        = "spoticli-db-sg"
  description = "Allow MySQL access from ECS tasks"
  vpc_id      = data.aws_vpc.default.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "spoticli-db-sg"
  }
}

resource "aws_security_group_rule" "db_from_tasks" {
  type                     = "ingress"
  description              = "MySQL from ECS tasks"
  from_port                = var.db_port
  to_port                  = var.db_port
  protocol                 = "tcp"
  security_group_id        = aws_security_group.db.id
  source_security_group_id = aws_security_group.fargate_tasks.id
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

# Database subnet group and instance
resource "aws_db_subnet_group" "default" {
  name       = "spoticli-db-subnets"
  subnet_ids = data.aws_subnets.default.ids

  tags = {
    Name = "spoticli-db-subnets"
  }
}

resource "aws_db_instance" "mysql" {
  identifier              = "spoticli-db"
  allocated_storage       = 20
  engine                  = "mysql"
  engine_version          = "8.0"
  instance_class          = "db.t3.micro"
  db_name                 = var.db_name
  username                = var.db_username
  password                = var.db_password
  port                    = var.db_port
  skip_final_snapshot     = true
  apply_immediately       = true
  backup_retention_period = 0
  publicly_accessible     = true
  db_subnet_group_name    = aws_db_subnet_group.default.name
  vpc_security_group_ids  = [aws_security_group.db.id]

  tags = {
    Name = "spoticli-mysql"
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
      },
      {
        name  = "DB_HOST"
        value = aws_db_instance.mysql.address
      },
      {
        name  = "DB_PORT"
        value = tostring(var.db_port)
      },
      {
        name  = "DB_USERNAME"
        value = var.db_username
      },
      {
        name  = "DB_PASSWORD"
        value = var.db_password
      },
      {
        name  = "STREAM_SEGMENT_SIZE"
        value = tostring(var.stream_segment_size)
      },
      {
        name  = "FRAME_CLUSTER_SIZE"
        value = tostring(var.frame_cluster_size)
      },
      {
        name  = "AWS_REGION"
        value = var.aws_region
      },
      {
        name  = "AWS_ENDPOINT_URL"
        value = var.aws_endpoint_url
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

output "db_endpoint" {
  description = "RDS endpoint for the application database"
  value       = aws_db_instance.mysql.address
}
