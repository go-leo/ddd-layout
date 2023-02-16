#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

echo "--- require started ---"
#pwd
go get -u -d ./...
go get github.com/ClickHouse/clickhouse-go@v1.5.4
go get github.com/go-leo/leo/v2@feature/v2
go get gorm.io/driver/clickhouse@v0.3.3
go mod tidy
echo "--- require finished ---"
