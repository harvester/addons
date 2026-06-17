ROOT := $(realpath $(dir $(realpath $(firstword $(MAKEFILE_LIST)))))
DOCKER_BUILDKIT := 1
export DOCKER_BUILDKIT

ifdef CI
	BOLD  :=
	CYAN  :=
	RESET :=
else
	BOLD  := \033[1m
	CYAN  := \033[36m
	RESET := \033[0m
endif
BANNER = @printf "$(BOLD)$(CYAN)[target: $@]$(RESET)\n"

MK_HOST_ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
export MK_HOST_ARCH

MK_REPO             := github.com/harvester/addons
MK_REPO_ID          := $(shell printf '%s' "$(ROOT)$(MK_SYSTEM_ID)" | sha256sum | cut -c1-8)
MK_PROVIDER_VERSION := $(shell git describe --tags --always --dirty)
MK_CODECOV_TOKEN    ?=
MK_DOCKER_PROGRESS  ?= plain

MK_CODECOV_SECRET_ARG  := --secret id=codecov_token_$(MK_REPO_ID),env=MK_CODECOV_TOKEN --no-cache-filter=test
MK_GOLANGCI_LINT_IMAGE := golangci/golangci-lint:v2.8.0-alpine@sha256:1194f3bfcbaeeb92d8d159fdfbe2a79d18ec0a222d9d984b1438906bca416b51

MK_HELM_VERSION=v3.20.0
MK_HELM_SHA256_amd64=dbb4c8fc8e19d159d1a63dda8db655f9ffa4aac1b9a6b188b34a40957119b286
MK_HELM_SHA256_arm64=bfb14953295d5324d47ab55f3dfba6da28d46c848978c8fbf412d4271bdc29f1

DOCKER_BUILD := \
	docker build \
		--progress=$(MK_DOCKER_PROGRESS) \
		--build-arg REPO=$(MK_REPO) \
		--build-arg REPO_ID=$(MK_REPO_ID) \
		--build-arg HOST_ARCH=$(MK_HOST_ARCH) \
		--build-arg GOLANGCI_LINT_IMAGE=$(MK_GOLANGCI_LINT_IMAGE) \
		--build-arg ARCH=$(MK_HOST_ARCH) \
		--build-arg HELM_VERSION=${MK_HELM_VERSION} \
		--build-arg HELM_SHA256_amd64=${MK_HELM_SHA256_amd64} \
		--build-arg HELM_SHA256_arm64=${MK_HELM_SHA256_arm64} \
		-f $(ROOT)/Dockerfile $(ROOT)

.PHONY: default generate templates addons patch-charts clean
.DEFAULT: default

default: generate patch-charts
generate: templates addons

output:
	@mkdir -p ./output

templates: output
	$(BANNER)
	$(DOCKER_BUILD) --target $@-output --output type=local,dest=.

addons: output
	$(BANNER)
	$(DOCKER_BUILD) --target $@-output --output type=local,dest=.

patch-charts: output
	$(BANNER)
	$(DOCKER_BUILD) --target $@-output --output type=local,dest=.

clean:
	@rm -rf output
