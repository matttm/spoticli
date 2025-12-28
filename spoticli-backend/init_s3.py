import boto3
import os

s3_client = boto3.client(
    "s3",
    endpoint_url=f"http://localhost:4566",
    aws_access_key_id="test",
    aws_secret_access_key="test"
)

bucket_name = os.environ.get("SPOTICLI_TRACKS", "spoticli-tracks")
s3_client.create_bucket(Bucket=bucket_name)