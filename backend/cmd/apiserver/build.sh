#!/bin/bash
cd "$(dirname "$0")" || exit # Set current directory to where the script is.
cd ../../

docker buildx build \
  -t apiserver \
  --build-arg pat="$GITHUB_PAT" \
  -f cmd/apiserver/Dockerfile \
  --progress plain \
  --no-cache \
  .