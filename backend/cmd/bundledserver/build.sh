#!/bin/bash
cd "$(dirname "$0")" || exit # Set current directory to where the script is.
cd ../../../

docker buildx build \
  -t bundledserver \
  -f backend/cmd/bundledserver/Dockerfile \
  --secret id=ssh_key,src="$SSH_KEY_SOURCE" \
  --secret id=env_json,src="$ENV_JSON_SOURCE" \
  --progress plain \
  --no-cache \
  .