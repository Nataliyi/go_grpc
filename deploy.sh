#!/bin/sh

# # set project
# gcloud config set project clusterization-services

# # set location
# gcloud config set eventarc/location us-central1
# gcloud config set run/region us-central1

# # create custom service account
# gcloud beta iam service-accounts create cl-server

# # allow for custom service account
# gcloud projects add-iam-policy-binding clusterization-services \
#   --role=roles/servicemanagement.configEditor \
#   --member serviceAccount:cl-server@clusterization-services.iam.gserviceaccount.com

# gcloud projects add-iam-policy-binding clusterization-services \
#     --member "serviceAccount:cl-server@clusterization-services.iam.gserviceaccount.com" \
#     --role roles/servicemanagement.serviceController

# gcloud projects add-iam-policy-binding clusterization-services \
# --member "serviceAccount:cl-server@clusterization-services.iam.gserviceaccount.com" \
# --role "roles/iap.httpsResourceAccessor"

# protoc
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative ./protos/grpc.proto

# # build image
gcloud builds submit --tag gcr.io/clusterization-services/cl-server:v0.0.1

# # deploy container
gcloud beta run deploy cl-server \
    --image gcr.io/clusterization-services/cl-server:v0.0.1 \
    --platform=managed \
    --region=us-central1 \
    --allow-unauthenticated \
    --project=clusterization-services \
    --use-http2 \
    --cpu=2 \
    --memory=2G \
    --min-instances=1 \
    --set-env-vars=REDISHOST=10.185.158.92,REDISPORT=6379
    # --vpc-connector=cluster-cache-connector 

# # enable sevices
# gcloud services enable servicemanagement.googleapis.com
# gcloud services enable servicecontrol.googleapis.com
# gcloud services enable endpoints.googleapis.com

# check server
grpcurl \
    -insecure \
    -proto grpc.proto \
    -d '{"pid":44, "sid": 22}' \
    cl-server-now57mm4pa-uc.a.run.app:443 \
    ClusterizationAPI.StreamClasterization

grpcurl \
    -proto grpc.proto \
    -d '{"pid":44, "sid": 22}' \
    cl-server-now57mm4pa-uc.a.run.app:443 \
    ClusterizationAPI.UnaryClasterization
