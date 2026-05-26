#!/bin/bash
gcloud functions deploy get-device-info \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=nodejs24 \
  --source=. \
  --entry-point=handler \
  --trigger-http \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME,JWT_SECRET_KEY=$JWT_SECRET_KEY \
  --allow-unauthenticated
