#!/bin/bash
gcloud functions deploy send-email \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=go126 \
  --source=. \
  --entry-point=Handler \
  --trigger-topic=send-email-topic \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME,JWT_SECRET_KEY=$JWT_SECRET_KEY,RESEND_API_KEY=$RESEND_API_KEY \
  --allow-unauthenticated
