# Using Terraform with LocalStack (No Real AWS!)

This setup lets you practice Terraform commands locally without AWS costs or account access.

## Quick Start

### 1. Start LocalStack
```bash
docker compose up -d localstack
```

### 2. Practice Terraform Commands
```bash
# Initialize (download providers)
terraform init

# See what will be created
terraform plan

# Create resources in LocalStack
terraform apply

# View created resources
terraform show

# Destroy everything
terraform destroy
```

## How It Works

- **`overrides.tf`** redirects ALL AWS API calls to LocalStack (localhost:4566)
- **LocalStack** simulates AWS services (S3, RDS, ECS, IAM, CloudWatch, etc.)
- **Your Terraform code** stays the same - you're learning real Terraform patterns

## Switching to Real AWS

When ready to deploy to actual AWS:

```bash
# Disable LocalStack override
mv overrides.tf overrides.tf.disabled

# Now Terraform uses real AWS
terraform plan
terraform apply
```

## Verify LocalStack Resources

```bash
# List S3 buckets
aws --endpoint-url=http://localhost:4566 s3 ls

# List RDS instances
aws --endpoint-url=http://localhost:4566 rds describe-db-instances

# List ECS clusters
aws --endpoint-url=http://localhost:4566 ecs list-clusters
```

## Limitations

LocalStack free tier has some limitations:
- RDS creates logical instances (not actual MySQL servers)
- ECS/Fargate tasks are simulated
- Some advanced features require LocalStack Pro

But it's **perfect for learning Terraform syntax, state management, and workflows!**

## Troubleshooting

**Issue**: `terraform apply` fails with connection errors  
**Fix**: Ensure LocalStack is running: `docker ps | grep localstack`

**Issue**: Resources not showing up  
**Fix**: Check LocalStack logs: `docker logs spoticli-localstack`
