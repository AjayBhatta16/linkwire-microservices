#!/bin/bash
gcloud functions deploy post-click \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=nodejs24 \
  --source=. \
  --entry-point=handler \
  --trigger-topic=post-click-topic \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME,JWT_SECRET_KEY=$JWT_SECRET_KEY \
  --allow-unauthenticated
