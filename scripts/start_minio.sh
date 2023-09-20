#!/usr/bin/env bash

project_dir=$(git rev-parse --show-toplevel)
minio_data_dir=$project_dir/data/minio
mkdir -p $minio_data_dir
chmod +x $minio_data_dir

docker run -d \
-p 9000:9000 \
-p 9001:9001 \
--name minio \
-v $minio_data_dir:/data \
-e "MINIO_ROOT_USER=access_key" \
-e "MINIO_ROOT_PASSWORD=secret_key" \
minio/minio \
server /data --console-address ":9001"