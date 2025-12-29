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

variable "tracks_bucket_name" {
  description = "S3 bucket used for storing track objects"
  type        = string
  default     = "spoticli-tracks"
}

variable "db_name" {
  description = "Database name for MySQL"
  type        = string
  default     = "SPOTICLI_DB"
}

variable "db_username" {
  description = "Database username"
  type        = string
  default     = "ADMIN"
}

variable "db_password" {
  description = "Database password"
  type        = string
  default     = "ADMINADMIN"
  sensitive   = true

  validation {
    condition     = length(var.db_password) >= 8
    error_message = "db_password must be at least 8 characters to satisfy RDS requirements."
  }
}

variable "db_port" {
  description = "MySQL port"
  type        = number
  default     = 3306
}

variable "stream_segment_size" {
  description = "STREAM_SEGMENT_SIZE environment variable"
  type        = number
  default     = 1000000
}

variable "frame_cluster_size" {
  description = "FRAME_CLUSTER_SIZE environment variable"
  type        = number
  default     = 600
}

variable "aws_endpoint_url" {
  description = "Optional AWS endpoint override (e.g., http://localstack:4566)"
  type        = string
  default     = ""
}
