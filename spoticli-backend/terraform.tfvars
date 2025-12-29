# Variable Values (referenced from variables.tf)
aws_region     = "us-east-1"
app_port       = 4200
fargate_cpu    = 256   # 0.25 vCPU
fargate_memory = 512   # 512 MB

# Full image reference (ECR, Docker Hub, etc.)
# Example for ECR: "123456789012.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest"
# Example for Docker Hub: "dockerhub-user/spoticli-backend:latest"
app_image = "123456789012.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest"

tracks_bucket_name  = "spoticli-tracks"
db_name             = "SPOTICLI_DB"
db_username         = "ADMIN"
db_password         = "ADMINADMIN" # Must be at least 8 characters for RDS
db_port             = 3306
stream_segment_size = 1000000
frame_cluster_size  = 600

# Set to http://localstack:4566 when using LocalStack (or leave empty for real AWS)
aws_endpoint_url = ""