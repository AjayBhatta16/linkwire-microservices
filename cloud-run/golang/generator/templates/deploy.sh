#!/bin/bash
gcloud functions deploy {{FUNCTION_NAME}} \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=go126 \
  --source=. \
  --entry-point=Handler \
  --trigger-http \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME \
  --allow-unauthenticated