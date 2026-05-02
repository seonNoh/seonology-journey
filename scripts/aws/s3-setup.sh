#!/usr/bin/env bash
# S3 buckets for seonology-journey: media + backup.
# Idempotent.
set -euo pipefail

PROFILE="${AWS_PROFILE:-seonology}"
REGION="${AWS_REGION:-ap-northeast-1}"

BUCKETS=(
  "seonology-journey-media"
  "seonology-journey-backup"
)

aws_cli() { aws --profile "$PROFILE" --region "$REGION" "$@"; }

for B in "${BUCKETS[@]}"; do
  if aws_cli s3api head-bucket --bucket "$B" 2>/dev/null; then
    echo "[skip] $B exists"
  else
    echo "[create] $B"
    aws_cli s3api create-bucket \
      --bucket "$B" \
      --create-bucket-configuration "LocationConstraint=$REGION"
  fi

  echo "[block-public] $B"
  aws_cli s3api put-public-access-block \
    --bucket "$B" \
    --public-access-block-configuration \
      "BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true"

  echo "[encryption] $B"
  aws_cli s3api put-bucket-encryption \
    --bucket "$B" \
    --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'

  echo "[versioning] $B"
  aws_cli s3api put-bucket-versioning \
    --bucket "$B" \
    --versioning-configuration Status=Enabled
done

# media バケット: CORS for direct PUT from browser
echo "[cors] seonology-journey-media"
aws_cli s3api put-bucket-cors --bucket seonology-journey-media --cors-configuration '{
  "CORSRules": [
    {
      "AllowedOrigins": ["https://journey.seonology.com"],
      "AllowedMethods": ["GET","PUT","POST","HEAD"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag"],
      "MaxAgeSeconds": 3000
    }
  ]
}'

# media: lifecycle - thumbnails/* に IA 移行 (90日) など
echo "[lifecycle] seonology-journey-media"
aws_cli s3api put-bucket-lifecycle-configuration --bucket seonology-journey-media --lifecycle-configuration '{
  "Rules": [
    {
      "ID": "noncurrent-cleanup",
      "Status": "Enabled",
      "Filter": {"Prefix": ""},
      "NoncurrentVersionExpiration": {"NoncurrentDays": 30},
      "AbortIncompleteMultipartUpload": {"DaysAfterInitiation": 7}
    }
  ]
}'

# backup: lifecycle - 30日 IA, 90日 Glacier
echo "[lifecycle] seonology-journey-backup"
aws_cli s3api put-bucket-lifecycle-configuration --bucket seonology-journey-backup --lifecycle-configuration '{
  "Rules": [
    {
      "ID": "tiering",
      "Status": "Enabled",
      "Filter": {"Prefix": ""},
      "Transitions": [
        {"Days": 30, "StorageClass": "STANDARD_IA"},
        {"Days": 90, "StorageClass": "GLACIER"}
      ],
      "AbortIncompleteMultipartUpload": {"DaysAfterInitiation": 7}
    }
  ]
}'

echo "buckets ready."
