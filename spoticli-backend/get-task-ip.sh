#!/bin/bash
# Script to get the public IP of the running Fargate task

set -e

AWS_REGION="${AWS_REGION:-us-east-1}"
AWS_PROFILE="${AWS_PROFILE:-matttm}"
CLUSTER_NAME="spoticli-cluster"
SERVICE_NAME="spoticli-backend-service"

echo "Getting task IP for $SERVICE_NAME in cluster $CLUSTER_NAME..."

TASK_ARN=$(aws ecs list-tasks \
  --cluster $CLUSTER_NAME \
  --service-name $SERVICE_NAME \
  --profile $AWS_PROFILE \
  --region $AWS_REGION \
  --query 'taskArns[0]' \
  --output text)

if [ "$TASK_ARN" == "None" ] || [ -z "$TASK_ARN" ]; then
  echo "No tasks found for service $SERVICE_NAME"
  exit 1
fi

echo "Task ARN: $TASK_ARN"

ENI_ID=$(aws ecs describe-tasks \
  --cluster $CLUSTER_NAME \
  --tasks $TASK_ARN \
  --profile $AWS_PROFILE \
  --region $AWS_REGION \
  --query 'tasks[0].attachments[0].details[?name==`networkInterfaceId`].value' \
  --output text)

echo "Network Interface ID: $ENI_ID"

PUBLIC_IP=$(aws ec2 describe-network-interfaces \
  --network-interface-ids $ENI_ID \
  --profile $AWS_PROFILE \
  --region $AWS_REGION \
  --query 'NetworkInterfaces[0].Association.PublicIp' \
  --output text)

echo ""
echo "========================================="
echo "Public IP: $PUBLIC_IP"
echo "Application URL: http://$PUBLIC_IP:4200"
echo "========================================="
