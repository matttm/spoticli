#!/bin/bash
set -e

# Deployment script for Spoticli Backend to AWS Fargate

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
AWS_REGION="${AWS_REGION:-us-east-1}"
AWS_PROFILE="${AWS_PROFILE:-matttm}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo -e "${GREEN}=== Spoticli Backend Deployment ===${NC}"
echo ""

# Get AWS Account ID
echo -e "${YELLOW}Getting AWS Account ID...${NC}"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --profile $AWS_PROFILE --query Account --output text)
echo -e "${GREEN}AWS Account ID: $AWS_ACCOUNT_ID${NC}"

# Get ECR Repository URL
ECR_REPO_NAME="spoticli-backend"
ECR_URL="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO_NAME}"

echo -e "${YELLOW}ECR Repository: $ECR_URL${NC}"
echo ""

# Login to ECR
echo -e "${YELLOW}Logging in to ECR...${NC}"
aws ecr get-login-password --region $AWS_REGION --profile $AWS_PROFILE | \
    docker login --username AWS --password-stdin $ECR_URL

# Build Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
docker build -t $ECR_REPO_NAME:$IMAGE_TAG -f Dockerfile.app .

# Tag image for ECR
echo -e "${YELLOW}Tagging image for ECR...${NC}"
docker tag $ECR_REPO_NAME:$IMAGE_TAG $ECR_URL:$IMAGE_TAG

# Push image to ECR
echo -e "${YELLOW}Pushing image to ECR...${NC}"
docker push $ECR_URL:$IMAGE_TAG

echo ""
echo -e "${GREEN}=== Deployment Complete ===${NC}"
echo -e "${GREEN}Image: $ECR_URL:$IMAGE_TAG${NC}"
echo ""
echo -e "${YELLOW}To update the ECS service, run:${NC}"
echo -e "aws ecs update-service --cluster spoticli-cluster --service spoticli-backend-service --force-new-deployment --profile $AWS_PROFILE --region $AWS_REGION"
