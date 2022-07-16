#!/bin/bash
set -eo pipefail

WORKDIR="$(dirname $(realpath "$0"))/../"
ROOT="$WORKDIR../"

cp "$WORKDIR/docker/iploc-server/Dockerfile" $ROOT
cd $ROOT

TARGET_IMAGE="iploc-server"
IMAGE_TAG="1.0.0"

docker build --tag "$TARGET_IMAGE:$IMAGE_TAG" . || exit 1

rm -f "$ROOT/Dockerfile"
