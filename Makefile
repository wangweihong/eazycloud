# .DEFAULT_GOAL为makefile自带变量, 用于设置默认目标
# https://www.gnu.org/software/make/manual/html_node/Special-Variables.html
.DEFAULT_GOAL := all

# Build options
ROOT_PACKAGE=github.com/wangweihong/eazycloud

.PHONY: all
all: tidy format lint cover build

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk
include scripts/make-rules/dependencies.mk
include scripts/make-rules/swagger.mk

# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="iam-apiserver iam-authz-server"
  IMAGES           Backend images to make. Default is all of cmd starting with iam-.
                   This option is available when using: make image/image.multiarch/push/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-authz-server"
  REGISTRY_PREFIX  Docker registry prefix. Default is marmotedu.
                   Example: make push REGISTRY_PREFIX=ccr.ccs.tencentyun.com/marmotedu VERSION=v1.6.2
  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multiarch/image.multiarch/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-pump" PLATFORMS="linux_amd64 linux_arm64"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch


## deploy: Deploy updated components to development env.
#.PHONY: deploy
#deploy:
#	@$(MAKE) deploy.run

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint:
	@$(MAKE) go.lint

## test: Run unit test.
#.PHONY: test
#test:
#	@$(MAKE) go.test

## cover: Run unit test and get test coverage.
#.PHONY: cover
#cover:
#	@$(MAKE) go.test.cover

## release: Release
#.PHONY: release
#release:
#	@$(MAKE) release.run

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .
	@$(GO) mod edit -fmt


## swagger: Generate swagger document.
#.PHONY: swagger
#swagger:
#	@$(MAKE) swagger.run

## serve-swagger: Serve swagger spec and docs.
#.PHONY: swagger.serve
#serve-swagger:
#	@$(MAKE) swagger.serve

## swagger-example: Generate example swagger and serve.
.PHONY: swagger.example
swagger-example:
	@$(MAKE) swagger.example
	@$(MAKE) swagger.example.serve

## dependencies: Install necessary dependencies.
.PHONY: dependencies
dependencies:
	@$(MAKE) dependencies.run

## tools: Install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

## check-updates: Check outdated dependencies of the go projects.
.PHONY: check-updates
check-updates:
	@$(MAKE) go.updates

## tidy: Go mod tidy
.PHONY: tidy
tidy:
	@$(GO) mod tidy

## deecopy-gen-example: Run deepcopy-gen example
.PHONY: deecopy-gen-example
deecopy-gen-example: tools.verify.deepcopy-gen
	@deepcopy-gen --input-dirs=./tools/deepcopy-gen/example --output-base=../

## help: Show this help info.
# 这里会提取target上一行的\#\#注释并生成到Makefile help文档中
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"
