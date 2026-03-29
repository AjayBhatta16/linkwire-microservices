gcloud functions deploy login \
  --gen2 \
  --region=$GCP_REGION \
  --runtime=go122 \
  --source=. \
  --entry-point=ProcessRequest \
  --trigger-http \
  --set-env-vars PROJECT_ID=$GCP_PROJECT_NAME \
  --allow-unauthenticated
