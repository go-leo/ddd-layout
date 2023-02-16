#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

echo "--- protoc generate start ---"
protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=module=media-adm \
  --go-grpc_out=. \
  --go-grpc_opt=module=media-adm \
  "$1"
echo "--- protoc generate end ---"
