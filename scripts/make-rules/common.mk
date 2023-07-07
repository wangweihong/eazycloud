SHELL := /bin/bash

# include the common make file
# MAKEFILE_LIST: makefile自带的环境变量，包含所有的makefile文件
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif
ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match='v*')
endif
# Check if the tree is dirty.  default to dirty
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

# Minimum test coverage
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 60
endif

# The OS must be linux when building docker images
PLATFORMS ?= linux_amd64 linux_arm64
# The OS can be linux/windows/darwin when building binaries
# PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

# Set a specific PLATFORM
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	# Use linux as the default OS when building images
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# Linux command settings
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
XARGS := xargs --no-run-if-empty

# Makefile settings
ifndef V
MAKEFLAGS += --no-print-directory
endif

ifeq ($(origin CHANGE_HOOK_LINE_SPERATOR), undefined)
	# 保证脚本换行符为\n,CRLF-->LF
	#CHANGE_HOOK_LINE_SPERATOR = $(shell dos2unix ./scripts/githooks/* )
	CHANGE_HOOK_LINE_SPERATOR = $(shell find ./scripts/githooks -type f -exec sh -c 'tr -d "\r" < "$0" > "$0.tmp" && mv "$0.tmp" "$0"' {} \; )
	# 保证脚本可执行
	MAKE_HOOK_EXECUTABLE:= $(shell chmod +x ./scripts/githooks/*)
    # Copy githook scripts when execute makefile
    # 采取这种方式, 可以实现git hook的统一和强制. 当执行Make任意规则时,强制进行拷贝。因此不需要单独的规则来拷贝
    COPY_GITHOOK:=$(shell cp -f ./scripts/githooks/* .git/hooks/)
endif

#endif
# Specify components which need certificate
#ifeq ($(origin CERTIFICATES),undefined)
#CERTIFICATES=iam-apiserver iam-authz-server admin
#endif

# Specify tools severity, include: BLOCKER_TOOLS, CRITICAL_TOOLS, TRIVIAL_TOOLS.
# Missing BLOCKER_TOOLS can cause the CI flow execution failed, i.e. `make all` failed.
# Missing CRITICAL_TOOLS can lead to some necessary operations failed. i.e. `make release` failed.
# TRIVIAL_TOOLS are Optional tools, missing these tool have no affect.
#BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint addlicense goimports codegen
BLOCKER_TOOLS ?= gsemver golines go-junit-report golangci-lint goimports codegen deepcopy-gen
#CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog github-release go-mod-outdated protoc-gen-go cfssl
CRITICAL_TOOLS ?= swagger mockgen gotests git-chglog  go-mod-outdated protoc-gen-go go-gitlint
#TRIVIAL_TOOLS ?= depth go-callvis gothanks richgo rts kube-score
TRIVIAL_TOOLS ?= depth go-callvis  richgo rts kube-score

COMMA := ,
EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
