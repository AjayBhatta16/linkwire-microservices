#!/bin/bash
gcloud functions deploy process-link \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=go126 \
  --source=. \
  --entry-point=Handler \
  --trigger-topic=process-link-topic \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME \
  --allow-unauthenticated
