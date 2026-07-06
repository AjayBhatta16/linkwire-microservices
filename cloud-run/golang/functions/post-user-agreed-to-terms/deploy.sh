#!/bin/bash
gcloud functions deploy post-user-agreed-to-terms \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=go126 \
  --source=. \
  --entry-point=Handler \
  --trigger-http \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME,JWT_SECRET_KEY=$JWT_SECRET_KEY \
  --allow-unauthenticated
