# note: call scripts from /scripts

GIT_VERSION=$(shell git tag | grep v | sort -r --version-sort | head -n1)
PROJECT_PATH:=$(shell pwd)

.PHONY: wire_gen
wire_gen:
	@sh scripts/shell/wire.sh

.PHONY: go_gen
go_gen:
	@sh scripts/shell/gen.sh


export protoFiles=$(shell listfile -ext=.proto)
.PHONY: protoc_gen
protoc_gen:
	@sh scripts/shell/protoc.sh $(protoFiles)

.PHONY: gen_ddd
gen_ddd:
	@echo "--- generate ddd app ---"
	genddd -module media-adm  -app $(app)
	@echo "--- generate ddd app end ---"

.PHONY: gen_all
gen_all:
	@echo "--- generate code start ---"
	@$(MAKE) wire_gen
	@$(MAKE) go_gen
	@$(MAKE) protoc_gen
	@echo "--- generate code end ---"

.PHONY: lint
lint:
	@sh scripts/shell/lint.sh

.PHONY: test
test:
	@echo "--- go test start ---"
	go test -test.bench=".*" -count=1 -v ./...
	@echo "--- go test end ---"

.PHONY: build
build:
	@echo "--- go build start ---"
	go build -o bin/api -a -ldflags "-w -s -X main.Version=$(GIT_VERSION)" -tags=jsoniter ./cmd/api/.
	@echo "--- go build end ---"

.PHONY: tools
tools:
	@sh scripts/shell/tools.sh


.PHONY: format
format:
	@sh scripts/shell/format.sh

.PHONY: require
require:
	@sh scripts/shell/require.sh

.PHONY: hooks
hooks:
	@cp ${PROJECT_PATH}/.gitx/hooks/* ${PROJECT_PATH}/.git/hooks/
	@chmod +x ${PROJECT_PATH}/.git/hooks/*

