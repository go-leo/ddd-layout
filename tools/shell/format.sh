#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

echo "--- go format start ---"
gofumpt -w ...
echo "--- go format end ---"

