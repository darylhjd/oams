#!/bin/bash
cd "$(dirname "$0")" || exit # Set current directory to where the script is.
cd ../../

docker buildx build \
  -t apiserver \
  -f cmd/apiserver/Dockerfile \
  --ssh default="$SSH_KEY_SOURCE" \
  --progress plain \
  --no-cache \
  .