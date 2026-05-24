#!/bin/bash
gcloud functions deploy {{FUNCTION_NAME}} \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=nodejs24 \
  --source=. \
  --entry-point=handler \
  --trigger-topic={{FUNCTION_NAME}}-topic \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME,JWT_SECRET_KEY=$JWT_SECRET_KEY \
  --allow-unauthenticated