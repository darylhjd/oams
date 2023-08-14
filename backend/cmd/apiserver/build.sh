#!/bin/bash
cd "$(dirname "$0")" || exit # Set current directory to where the script is.
cd ../../

docker buildx build \
  -t apiserver \
  -f cmd/apiserver/Dockerfile \
  --secret id=ssh_key,src="$SSH_KEY_SOURCE" \
  .