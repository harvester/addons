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

MK_SYSTEM_ID := $(strip $(shell \
		if [ -s /etc/machine-id ]; then \
				cat /etc/machine-id 2>/dev/null; \
		elif command -v hostname >/dev/null 2>&1; then \
				hostname 2>/dev/null; \
		else \
				echo -n "unknown"; \
		fi))

MK_REPO             := github.com/harvester/addons
MK_REPO_ID          := $(shell printf '%s' "$(ROOT)$(MK_SYSTEM_ID)" | sha256sum | cut -c1-8)
MK_PROVIDER_VERSION := $(shell git describe --tags --always --dirty)
MK_CODECOV_TOKEN    ?=
MK_DOCKER_PROGRESS  ?= plain

MK_GOLANGCI_LINT_IMAGE := golangci/golangci-lint:v2.8.0-alpine@sha256:1194f3bfcbaeeb92d8d159fdfbe2a79d18ec0a222d9d984b1438906bca416b51

DOCKER_BUILD := \
	docker build \
		--progress=$(MK_DOCKER_PROGRESS) \
		--build-arg MK_REPO=$(MK_REPO) \
		--build-arg MK_REPO_ID=$(MK_REPO_ID) \
		--build-arg MK_HOST_ARCH=$(MK_HOST_ARCH) \
		--build-arg MK_GOLANGCI_LINT_IMAGE=$(MK_GOLANGCI_LINT_IMAGE) \
		-f $(ROOT)/Dockerfile $(ROOT)

.PHONY: generate test-chart-patch
.DEFAULT_GOAL := default
default: generate test-chart-patch

# for github workflow usage
ci: generate test-chart-patch

# ---- Directories ----
$(ROOT)/bin:
	@mkdir -p $@

# ---- generate Addons ----
generate: $(ROOT)/bin
	$(BANNER)
	$(DOCKER_BUILD) --target generate-output --output type=local,dest=.

# ---- test chart patches ----
test-chart-patch: $(ROOT)/bin
	$(BANNER)
	$(DOCKER_BUILD) --target test-chart-patch-output --output type=local,dest=.
