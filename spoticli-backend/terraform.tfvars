# Copy this file to terraform.tfvars and fill in your values
# cp terraform.tfvars.example terraform.tfvars

aws_region     = "us-east-1"
app_port       = 4200
fargate_cpu    = 256   # 0.25 vCPU
fargate_memory = 512   # 512 MB
# Full image reference (ECR, Docker Hub, etc.)
# Example for ECR: "123456789012.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest"
# Example for Docker Hub: "dockerhub-user/spoticli-backend:latest"
app_image      = "123456789012.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest"
