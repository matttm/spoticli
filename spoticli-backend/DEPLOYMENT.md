# Spoticli Backend - AWS Fargate Deployment Guide

This Terraform configuration deploys the Spoticli backend to AWS using **Fargate** (serverless containers) instead of EC2 instances.

## Architecture

- **AWS Fargate**: Runs your Go application in Docker containers (serverless, no EC2 management)
- **Public IP**: Fargate task gets a public IP for direct access
- **ECR (Elastic Container Registry)**: Stores Docker images
- **ECS (Elastic Container Service)**: Orchestrates Fargate tasks
- **S3**: Two buckets - one for music storage, one for Terraform state
- **VPC**: Custom VPC with 2 public subnets across availability zones
- **CloudWatch Logs**: Captures application logs

## Cost Estimation (Monthly)

### Fargate Configuration (256 CPU, 512 MB Memory)
- **Fargate compute**: ~$8-10/month (running 24/7)
- **S3 storage**: $0.023 per GB + request costs
- **CloudWatch Logs**: ~$0.50/GB ingested
- **Data transfer**: First 100GB outbound free, then $0.09/GB
- **ECR storage**: $0.10 per GB/month

**Estimated Total: ~$10-15/month**

### Ways to Reduce Costs:
1. **Use AWS Free Tier** benefits (first 12 months)
2. **Scale down** when not in use: `aws ecs update-service --desired-count 0`
3. Reduce Fargate CPU/memory if possible
4. Use CloudWatch Logs retention policies

## Prerequisites

1. **AWS CLI** configured with credentials
2. **Docker** installed
3. **Terraform** installed (>= 1.0)
4. **AWS Account** with appropriate permissions

## Setup Instructions

### Step 1: Initial AWS Setup

Create the S3 bucket and DynamoDB table for Terraform state (one-time setup):

```bash
# Create S3 bucket for Terraform state
aws s3api create-bucket \
  --bucket spoticli-terraform-state \
  --region us-east-1 \
  --profile matttm

# Enable versioning
aws s3api put-bucket-versioning \
  --bucket spoticli-terraform-state \
  --versioning-configuration Status=Enabled \
  --profile matttm

# Create DynamoDB table for state locking
aws dynamodb create-table \
  --table-name spoticli-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1 \
  --profile matttm
```

### Step 2: Configure Terraform Variables

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars if needed (defaults should work)
```

### Step 3: Initialize Terraform

```bash
terraform init
```

### Step 4: Deploy Infrastructure

```bash
# Review the plan
terraform plan

# Apply the configuration
terraform apply
```

This will create:
- VPC with networking
- ECR repository
- ECS cluster
- IAM roles
- Security groups
- Application Load Balancer
- S3 buckets
- CloudWatch log group

**Note:** The ECS service will initially fail to start because there's no Docker image yet.

### Step 5: Build and Deploy Your Application

```bash
# Build and push Docker image to ECR
./deploy.sh

# Or manually:
# Get ECR login
aws ecr get-login-password --region us-east-1 --profile matttm | \
  docker login --username AWS --password-stdin \
  $(aws sts get-caller-identity --profile matttm --query Account --output text).dkr.ecr.us-east-1.amazonaws.com

# Build image
docker build -t spoticli-backend:latest -f Dockerfile.app .

# Tag and push
ACCOUNT_ID=$(aws sts get-caller-identity --profile matttm --query Account --output text)
docker tag spoticli-backend:latest $ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest
docker push $ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest

# Force ECS to deploy the new image
aws ecs update-service \
  --cluster spoticli-cluster \
  --service spoticli-backend-service \
  --force-new-deployment \
  --profile matttm \
  --region us-east-1
```

### Step 6: Access Your Application

Get the public IP of your running Fargate task:
```bash
# Using the helper command
terraform output -raw get_task_ip_command | bash

# Or manually
aws ecs list-tasks \
  --cluster spoticli-cluster \
  --service-name spoticli-backend-service \
  --profile matttm \
  --region us-east-1 \
  --query 'taskArns[0]' \
  --output text | xargs -I {} \
aws ecs describe-tasks \
  --cluster spoticli-cluster \
  --tasks {} \
  --profile matttm \
  --region us-east-1 \
  --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' \
  --output text | xargs -I {} \
aws ec2 describe-network-interfaces \
  --network-interface-ids {} \
  --profile matttm \
  --region us-east-1 \
  --query 'NetworkInterfaces[0].Association.PublicIp' \
  --output text
```

Access your application at `http://<TASK_PUBLIC_IP>:4200`

**Note:** The public IP will change if the task restarts. For a static endpoint, you could add an Application Load Balancer or use a service discovery solution.

## Useful Commands

### View Logs
```bash
# Stream logs
aws logs tail /ecs/spoticli-backend --follow --profile matttm --region us-east-1

# View specific task logs
aws ecs describe-tasks \
  --cluster spoticli-cluster \
  --tasks $(aws ecs list-tasks --cluster spoticli-cluster --service spoticli-backend-service --profile matttm --region us-east-1 --query 'taskArns[0]' --output text) \
  --profile matttm \
  --region us-east-1
```

### Scale Service
```bash
# Scale to 2 tasks
aws ecs update-service \
  --cluster spoticli-cluster \
  --service spoticli-backend-service \
  --desired-count 2 \
  --profile matttm \
  --region us-east-1

# Scale to 0 (stop service to save costs)
aws ecs update-service \
  --cluster spoticli-cluster \
  --service spoticli-backend-service \
  --desired-count 0 \
  --profile matttm \
  --region us-east-1
```

### Redeploy Application
```bash
# After pushing new image
aws ecs update-service \
  --cluster spoticli-cluster \
  --service spoticli-backend-service \
  --force-new-deployment \
  --profile matttm \
  --region us-east-1
```

### Check Service Status
```bash
aws ecs describe-services \
  --cluster spoticli-cluster \
  --services spoticli-backend-service \
  --profile matttm \
  --region us-east-1
```

## Outputs

After `terraform apply`, you'll see:

- `alb_url`: Load balancer URL to access your application
- `ecr_repository_url`: ECR repository for Docker images
- `music_bucket_name`: S3 bucket name for music storage
- `ecs_cluster_name`: ECS cluster name
- `ecs_service_name`: ECS service name
- `cloudwatch_log_group`: Log group for viewing application logs

## Environment Variables

The application receives these environment variables:
- `MUSIC_BUCKET`: Name of the S3 bucket for music storage
- `PORT`: Application port (4200)

Access the music bucket from your Go code using the AWS SDK.

## Cleanup

To destroy all resources:

```bash
# First, delete all images from ECR
aws ecr batch-delete-image \
  --repository-name spoticli-backend \
  --image-ids imageTag=latest \
  --profile matttm \
  --region us-east-1

# Then destroy infrastructure
terraform destroy
```

**Note:** The Terraform state bucket has `prevent_destroy = true`. You'll need to manually delete it if needed.

## Troubleshooting

### Service won't start
- Check ECS console for task failure reasons
- View logs: `aws logs tail /ecs/spoticli-backend --follow`
- Ensure Docker image exists in ECR
- Verify health check endpoint returns 200

### Can't connect to task
- Check security groups allow port 4200 from your IP
- Verify task is running: `aws ecs describe-services --cluster spoticli-cluster --services spoticli-backend-service`
- Ensure task has public IP assigned

### Application can't access S3
- Verify IAM task role has S3 permissions
- Check bucket name in environment variables
- Ensure tasks have internet access (public subnets with IGW)

## Next Steps

1. **Set up CI/CD**: Use GitHub Actions or AWS CodePipeline for automated deployments
2. **Add Load Balancer**: For production, add an ALB for a static endpoint and HTTPS support
3. **Custom Domain**: Use Route 53 to map a domain name
4. **Database**: Set up RDS for persistent storage
5. **Monitoring**: Configure CloudWatch alarms for health monitoring
6. **Auto-scaling**: Add ECS auto-scaling based on CPU/memory metrics
