#!/usr/bin/env bash
# DynamoDB tables for seonology-journey.
# Idempotent: 既存テーブルは skip.
#
# 前提: AWS_PROFILE=seonology, region=ap-northeast-1
set -euo pipefail

PROFILE="${AWS_PROFILE:-seonology}"
REGION="${AWS_REGION:-ap-northeast-1}"
PREFIX="seonology-journey-"

aws_cli() {
  aws --profile "$PROFILE" --region "$REGION" "$@"
}

table_exists() {
  aws_cli dynamodb describe-table --table-name "$1" >/dev/null 2>&1
}

create_table() {
  local name="$1"
  local args=("$@")
  args=("${args[@]:1}")
  local full="${PREFIX}${name}"
  if table_exists "$full"; then
    echo "[skip] $full"
    return 0
  fi
  echo "[create] $full"
  aws_cli dynamodb create-table \
    --table-name "$full" \
    --billing-mode PAY_PER_REQUEST \
    "${args[@]}"
  aws_cli dynamodb wait table-exists --table-name "$full"
  aws_cli dynamodb update-continuous-backups \
    --table-name "$full" \
    --point-in-time-recovery-specification PointInTimeRecoveryEnabled=true >/dev/null
  echo "[pitr-on] $full"
}

enable_ttl() {
  local name="$1"
  local attr="$2"
  local full="${PREFIX}${name}"
  aws_cli dynamodb update-time-to-live \
    --table-name "$full" \
    --time-to-live-specification "Enabled=true,AttributeName=${attr}" >/dev/null || true
  echo "[ttl] $full $attr"
}

# 1. trips: PK=USER#userId, SK=TRIP#tripId. GSI1 by tripId, GSI2 by status+startDate
create_table trips \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=S \
    AttributeName=GSI2PK,AttributeType=S \
    AttributeName=GSI2SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
    "[{\"IndexName\":\"GSI1\",\"KeySchema\":[{\"AttributeName\":\"GSI1PK\",\"KeyType\":\"HASH\"},{\"AttributeName\":\"GSI1SK\",\"KeyType\":\"RANGE\"}],\"Projection\":{\"ProjectionType\":\"ALL\"}},{\"IndexName\":\"GSI2\",\"KeySchema\":[{\"AttributeName\":\"GSI2PK\",\"KeyType\":\"HASH\"},{\"AttributeName\":\"GSI2SK\",\"KeyType\":\"RANGE\"}],\"Projection\":{\"ProjectionType\":\"ALL\"}}]"

# 2. days
create_table days \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 3. schedules
create_table schedules \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 4. meals (PK=DAY#dayId, SK=MEAL#type)
create_table meals \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 5. accommodations (PK=DAY#dayId, SK=ACCOMMODATION)
create_table accommodations \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 6. media (PK=TRIP#tripId, SK=MEDIA#takenAt#mediaId)
create_table media \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
    "[{\"IndexName\":\"GSI1\",\"KeySchema\":[{\"AttributeName\":\"GSI1PK\",\"KeyType\":\"HASH\"},{\"AttributeName\":\"GSI1SK\",\"KeyType\":\"RANGE\"}],\"Projection\":{\"ProjectionType\":\"ALL\"}}]"

# 7. expenses
create_table expenses \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 8. notes
create_table notes \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 9. companions
create_table companions \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
    "[{\"IndexName\":\"GSI1\",\"KeySchema\":[{\"AttributeName\":\"GSI1PK\",\"KeyType\":\"HASH\"},{\"AttributeName\":\"GSI1SK\",\"KeyType\":\"RANGE\"}],\"Projection\":{\"ProjectionType\":\"ALL\"}}]"

# 10. checklist
create_table checklist \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 11. reservations
create_table reservations \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 12. tags (PK=USER#userId, SK=TAG#tagId)
create_table tags \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 13. templates
create_table templates \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 14. favorite-places
create_table favorite-places \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE

# 15. shares (TTL: expiresAt)
create_table shares \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
  --global-secondary-indexes \
    "[{\"IndexName\":\"GSI1\",\"KeySchema\":[{\"AttributeName\":\"GSI1PK\",\"KeyType\":\"HASH\"},{\"AttributeName\":\"GSI1SK\",\"KeyType\":\"RANGE\"}],\"Projection\":{\"ProjectionType\":\"ALL\"}}]"
enable_ttl shares expiresAt

# 16. api-cache (TTL: expiresAt)
create_table api-cache \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH
enable_ttl api-cache expiresAt

echo "all tables ready."
