#!/usr/bin/env bash
# IAM user + policy for seonology-journey-back service.
# Generates an access key only on first creation; on re-run, skips.
set -euo pipefail

PROFILE="${AWS_PROFILE:-seonology}"
REGION="${AWS_REGION:-ap-northeast-1}"
USER="seonology-journey-back-svc"
POLICY_NAME="seonology-journey-back-policy"
ACCOUNT_ID=$(aws --profile "$PROFILE" sts get-caller-identity --query Account --output text)

aws_cli() { aws --profile "$PROFILE" --region "$REGION" "$@"; }

POLICY_DOC=$(cat <<JSON
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DynamoDBJourney",
      "Effect": "Allow",
      "Action": [
        "dynamodb:GetItem","dynamodb:PutItem","dynamodb:UpdateItem","dynamodb:DeleteItem",
        "dynamodb:Query","dynamodb:Scan","dynamodb:BatchGetItem","dynamodb:BatchWriteItem",
        "dynamodb:DescribeTable","dynamodb:TransactWriteItems","dynamodb:TransactGetItems"
      ],
      "Resource": [
        "arn:aws:dynamodb:${REGION}:${ACCOUNT_ID}:table/seonology-journey-*",
        "arn:aws:dynamodb:${REGION}:${ACCOUNT_ID}:table/seonology-journey-*/index/*"
      ]
    },
    {
      "Sid": "S3Media",
      "Effect": "Allow",
      "Action": [
        "s3:GetObject","s3:PutObject","s3:DeleteObject","s3:AbortMultipartUpload","s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::seonology-journey-media",
        "arn:aws:s3:::seonology-journey-media/*",
        "arn:aws:s3:::seonology-journey-backup",
        "arn:aws:s3:::seonology-journey-backup/*"
      ]
    }
  ]
}
JSON
)

# 1) policy upsert
POLICY_ARN="arn:aws:iam::${ACCOUNT_ID}:policy/${POLICY_NAME}"
if aws_cli iam get-policy --policy-arn "$POLICY_ARN" >/dev/null 2>&1; then
  echo "[policy] update version"
  aws_cli iam create-policy-version \
    --policy-arn "$POLICY_ARN" \
    --policy-document "$POLICY_DOC" \
    --set-as-default >/dev/null
  # 古いバージョンを最大4まで掃除
  for V in $(aws_cli iam list-policy-versions --policy-arn "$POLICY_ARN" --query 'Versions[?IsDefaultVersion==`false`].VersionId' --output text); do
    aws_cli iam delete-policy-version --policy-arn "$POLICY_ARN" --version-id "$V" || true
  done
else
  echo "[policy] create $POLICY_NAME"
  aws_cli iam create-policy \
    --policy-name "$POLICY_NAME" \
    --policy-document "$POLICY_DOC" >/dev/null
fi

# 2) user
if aws_cli iam get-user --user-name "$USER" >/dev/null 2>&1; then
  echo "[user] $USER exists"
else
  echo "[user] create $USER"
  aws_cli iam create-user --user-name "$USER" >/dev/null
fi

# 3) attach
aws_cli iam attach-user-policy --user-name "$USER" --policy-arn "$POLICY_ARN" || true
echo "[attached] $POLICY_ARN -> $USER"

# 4) access key (only if user has none)
KEY_COUNT=$(aws_cli iam list-access-keys --user-name "$USER" --query 'length(AccessKeyMetadata)' --output text)
if [[ "$KEY_COUNT" -eq 0 ]]; then
  echo "[key] creating access key (出力は一度のみ. Vault に登録すること)"
  aws_cli iam create-access-key --user-name "$USER" \
    --query '{ak:AccessKey.AccessKeyId,sk:AccessKey.SecretAccessKey}' \
    --output json
else
  echo "[key] $USER は既に access key を持つ ($KEY_COUNT). 必要なら手動で再発行."
fi

echo "iam ready."
