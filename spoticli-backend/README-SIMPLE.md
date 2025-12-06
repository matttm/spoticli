# Spoticli Backend - Simple AWS Deployment

This is a **simplified** Terraform configuration for beginners. It deploys your Go app to AWS Fargate with minimal complexity and assumes you already have a container image pushed to a registry (ECR, Docker Hub, etc).

## What This Creates

1. **S3 Bucket** - Stores your music files
2. **Fargate Service** - Runs your app in a container (no servers to manage!)
3. **Security Group** - Allows traffic to your app on port 4200
4. **IAM Roles** - Gives your app permission to access S3

**Uses your default VPC** - No custom networking needed!

## Cost: ~$10/month

- Fargate: ~$8-10/month (256 CPU, 512 MB memory)
- S3: Based on storage used
- Image registry: depends on where you host your image (ECR is ~$0.10/GB)

## Quick Start

### 1. Install Prerequisites

```bash
# Install Terraform
brew install terraform

# Install AWS CLI (if not already installed)
brew install awscli

# Configure AWS credentials (if not already done)
aws configure --profile matttm
```

### 2. Point to Your Image

Set the image you want Fargate to run (in `terraform.tfvars`):

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` and set `app_image` to the full image name (for example `123456789012.dkr.ecr.us-east-1.amazonaws.com/spoticli-backend:latest` or `dockerhub-user/spoticli-backend:latest`).

### 3. Build and Push Your App Image (if needed)

If you need to push your image to ECR, make sure the repository already exists and run:

```bash
# Build and push your Docker image
./deploy.sh
```

If you're using Docker Hub or another registry, push there instead and ensure `app_image` matches the pushed tag.

### 4. Deploy Infrastructure

With `app_image` set and the image available in your registry:

```bash
# Initialize Terraform
terraform init

# See what will be created
terraform plan

# Create everything!
terraform apply
```

Type `yes` when prompted.

### 5. Get Your App URL

```bash
# Get the public IP of your running app
./get-task-ip.sh
```

Visit `http://<IP>:4200` in your browser!

## Common Commands

**View logs:**
```bash
aws logs tail /ecs/spoticli-backend --follow --profile matttm
```

**Redeploy after code changes:**
```bash
./deploy.sh
```

**Stop the app (save money):**
```bash
aws ecs update-service \
  --cluster spoticli-cluster \
  --service spoticli-backend-service \
  --desired-count 0 \
  --profile matttm \
  --region us-east-1
```

**Start it again:**
```bash
aws ecs update-service \
  --cluster spoticli-cluster \
  --service spoticli-backend-service \
  --desired-count 1 \
  --profile matttm \
  --region us-east-1
```

**Delete everything:**
```bash
terraform destroy
```

## Simplified vs Full Production

This simplified setup removes:
- ❌ Custom VPC/networking
- ❌ Load balancer
- ❌ Remote state storage (S3 backend)
- ❌ Health checks
- ❌ Auto-scaling
- ❌ Multiple availability zones
- ❌ Container Insights
- ❌ Image lifecycle policies
- ❌ Versioning/encryption on S3

**This is perfect for:**
- Learning
- Development
- Simple side projects
- Low-traffic apps

**You'll need more for:**
- Production workloads
- High availability
- Team collaboration
- Compliance requirements

## Troubleshooting

**Can't connect to the app?**
- Make sure the task is running: `aws ecs describe-services --cluster spoticli-cluster --services spoticli-backend-service --profile matttm`
- Check logs: `aws logs tail /ecs/spoticli-backend --follow --profile matttm`

**Task keeps restarting?**
- Check that the image referenced in `app_image` exists and is publicly accessible to ECS (ECR repo permissions, Docker Hub, etc.)
- View container logs for errors

**"No default VPC" error?**
- Your AWS account doesn't have a default VPC. You'll need to create one or use a custom VPC setup.

## Next Steps

When you're ready to scale up:
1. Add a load balancer for a static URL
2. Enable auto-scaling based on traffic
3. Set up a custom domain with Route 53
4. Add RDS database for persistent data
5. Use remote state (S3 backend) for team collaboration
6. Add SSL/HTTPS support

## File Structure

- `main.tf` - All your infrastructure
- `Dockerfile.app` - Builds your Go app
- `deploy.sh` - Builds and pushes Docker image
- `get-task-ip.sh` - Gets your app's public IP
- `terraform.tfvars.example` - Configuration (optional)
