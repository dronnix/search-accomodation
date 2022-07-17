#!/bin/bash
set -eo pipefail

WORKDIR="$(dirname $(realpath "$0"))/../"
ROOT="$WORKDIR../"

cp "$WORKDIR/docker/iploc-data-importer/Dockerfile" $ROOT
cd $ROOT

TARGET_IMAGE="iploc-data-importer"
IMAGE_TAG="1.0.0"

docker build --tag "$TARGET_IMAGE:$IMAGE_TAG" . || exit 1

rm -f "$ROOT/Dockerfile"
