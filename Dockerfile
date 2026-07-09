# Provide a default value so the static linter is happy, it eliminates following warning:
# InvalidDefaultArgInFrom: Default value for ARG ${MK_GOLANGCI_LINT_IMAGE} results in empty or invalid base image name (line 2)
ARG MK_GOLANGCI_LINT_IMAGE=golangci/golangci-lint:latest
FROM ${MK_GOLANGCI_LINT_IMAGE} AS golangci-lint

FROM registry.suse.com/bci/golang:1.25.0 AS buildenv

ARG HELM_VERSION=v3.20.0
ARG HELM_SHA256_Linux_amd64=dbb4c8fc8e19d159d1a63dda8db655f9ffa4aac1b9a6b188b34a40957119b286
ARG HELM_SHA256_Linux_arm64=bfb14953295d5324d47ab55f3dfba6da28d46c848978c8fbf412d4271bdc29f1

RUN zypper -n rm container-suseconnect && \
    zypper -n install curl tar patch

# set up helm
ARG MK_HOST_ARCH=amd64
ARG MK_REPO_ID
ENV HELM_VERSION=${HELM_VERSION}
ENV HELM_TARBALL=helm-${HELM_VERSION}-linux-${MK_HOST_ARCH}.tar.gz
ENV HELM_URL=https://get.helm.sh/${HELM_TARBALL}

SHELL [ "/bin/bash", "-e", "-o", "pipefail", "-c" ]
RUN --mount=type=cache,target=/tmp/helm,id=helm-dl-${MK_REPO_ID} \
    <<EOF
#!/bin/bash

curl -sSLO --output-dir /tmp/helm ${HELM_URL}
HELM_SHA256=HELM_SHA256_Linux_${MK_HOST_ARCH}
echo "${!HELM_SHA256}  /tmp/helm/${HELM_TARBALL}"
echo "${!HELM_SHA256}  /tmp/helm/${HELM_TARBALL}" | sha256sum -c -
tar xvzf /tmp/helm/${HELM_TARBALL} --strip-components=1 -C /tmp/helm
mv /tmp/helm/helm /usr/bin/helm

EOF

COPY --from=golangci-lint /usr/bin/golangci-lint /usr/local/bin/golangci-lint

# ---- base ----
FROM buildenv AS base
ARG MK_REPO
ARG MK_REPO_ID
ENV DAPPER_SOURCE=/go/src/${MK_REPO}
WORKDIR /go/src/${MK_REPO}
# to exclude some files, add them in .dockerignore
COPY . .

# ---- generate ----
FROM base AS generate
ARG MK_REPO
ARG MK_REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${MK_REPO_ID} \
    --mount=type=cache,target=/go/src/${MK_REPO}/.cache/go-build,id=harvester-go-build-${MK_REPO_ID} \
    scripts/generate

# ---- test-chart-patch ----
FROM base AS test-chart-patch
ARG MK_REPO
ARG MK_REPO_ID
RUN  scripts/test-chart-patch

# ---- generate output ----
FROM scratch AS generate-output
ARG MK_REPO
COPY --from=generate /go/src/${MK_REPO}/bin/ /bin/

# ---- test-chart-patch output ----
FROM scratch AS test-chart-patch-output
ARG MK_REPO
COPY --from=test-chart-patch /go/src/${MK_REPO}/bin/ /bin/
