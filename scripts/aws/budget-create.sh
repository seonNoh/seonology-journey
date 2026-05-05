#!/usr/bin/env bash
set -euo pipefail

ACCOUNT_ID="756541174912"
PROFILE="seonology"
BUDGET_NAME="seonology-journey-monthly"
LIMIT="10"
EMAIL="seon@seonology.com"

AWS_PAGER="" AWS_PROFILE="$PROFILE" aws budgets create-budget \
  --account-id "$ACCOUNT_ID" \
  --budget "{
    \"BudgetName\": \"$BUDGET_NAME\",
    \"BudgetLimit\": {\"Amount\": \"$LIMIT\", \"Unit\": \"USD\"},
    \"TimeUnit\": \"MONTHLY\",
    \"BudgetType\": \"COST\",
    \"CostTypes\": {
      \"IncludeTax\": true,
      \"IncludeSubscription\": true,
      \"UseBlended\": false,
      \"IncludeRefund\": false,
      \"IncludeCredit\": false,
      \"IncludeUpfront\": true,
      \"IncludeRecurring\": true,
      \"IncludeOtherSubscription\": true,
      \"IncludeSupport\": true,
      \"IncludeDiscount\": true,
      \"UseAmortized\": false
    }
  }" \
  --notifications-with-subscribers "[
    {
      \"Notification\": {
        \"NotificationType\": \"ACTUAL\",
        \"ComparisonOperator\": \"GREATER_THAN\",
        \"Threshold\": 80,
        \"ThresholdType\": \"PERCENTAGE\"
      },
      \"Subscribers\": [{\"SubscriptionType\": \"EMAIL\", \"Address\": \"$EMAIL\"}]
    },
    {
      \"Notification\": {
        \"NotificationType\": \"ACTUAL\",
        \"ComparisonOperator\": \"GREATER_THAN\",
        \"Threshold\": 100,
        \"ThresholdType\": \"PERCENTAGE\"
      },
      \"Subscribers\": [{\"SubscriptionType\": \"EMAIL\", \"Address\": \"$EMAIL\"}]
    }
  ]"

echo "Budget '$BUDGET_NAME' created: \$$LIMIT/month, alerts at 80% and 100%"
