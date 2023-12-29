#!/bin/bash
cd "$(dirname "$0")" || exit # Set current directory to where the script is.
cd ../../../

docker buildx build \
  -t webserver \
  -f backend/cmd/webserver/Dockerfile \
  .