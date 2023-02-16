#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

echo "--- wire generate start ---"
wire ./...
echo "--- wire generate end ---"

