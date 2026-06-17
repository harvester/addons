ARG GOLANGCI_LINT_IMAGE
FROM ${GOLANGCI_LINT_IMAGE} AS golangci-lint

FROM golang:1.25.7-trixie AS buildenv

SHELL [ "/usr/bin/bash", "-e", "-o", "pipefail" ]

ARG ARCH
ARG HELM_VERSION
ENV HELM_VERSION=${HELM_VERSION}
ARG HELM_SHA256_amd64
ENV HELM_SHA256_amd64=${HELM_SHA256_amd64}
ARG HELM_SHA256_arm64
ENV HELM_SHA256_arm64=${HELM_SHA256_arm64}
ENV HELM_TARBALL=helm-${HELM_VERSION}-linux-${ARCH}.tar.gz
ENV HELM_URL=https://get.helm.sh/${HELM_TARBALL}
ENV HELM_SHA256=HELM_SHA256_${ARCH}

# hadolint ignore=DL3008
RUN --mount=type=cache,target=/var/lib/apt/lists <<EOF
#!/bin/bash -e -o pipefail

apt-get update -qq
apt-get install -y --no-install-recommends \
  curl \
  patch \
  tar

mkdir /tmp/helm
curl -sSLO --output-dir /tmp/helm ${HELM_URL}
echo "${!HELM_SHA256}  /tmp/helm/${HELM_TARBALL}" | sha256sum -c -
tar xvzf /tmp/helm/${HELM_TARBALL} --strip-components=1 -C /tmp/helm
mv /tmp/helm/helm /usr/bin/helm
EOF

# ---- base ----
FROM buildenv AS base
ARG REPO
ARG REPO_ID
WORKDIR /go/src/${REPO}
# to exclude some files, add them in .dockerignore
COPY . .

# ---- templates ----
FROM base AS templates
ARG REPO
ARG REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${REPO_ID} \
    --mount=type=cache,target=/go/src/${REPO}/.cache/go-build,id=harvester-go-build-${REPO_ID} \
    <<EOF
#!/bin/bash -e
mkdir -p "$(pwd)/output/templates"
go run . -generateTemplates -path "./output/templates"
EOF

# ---- addons ----
FROM base AS addons
ARG REPO
ARG REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${REPO_ID} \
    --mount=type=cache,target=/go/src/${REPO}/.cache/go-build,id=harvester-go-build-${REPO_ID} \
    <<EOF
#!/bin/bash -e
mkdir -p "$(pwd)/output/addons"
go run . -generateAddons -path "./output/addons"
EOF

# ---- patch-charts ----
FROM base AS patch-charts
ARG REPO
ARG REPO_ID
RUN --mount=type=cache,target=/go/pkg/mod,id=harvester-go-mod-${REPO_ID} \
    --mount=type=cache,target=/go/src/${REPO}/.cache/go-build,id=harvester-go-build-${REPO_ID} \
    scripts/patch-charts

# ---- templates output ----
FROM scratch AS templates-output
ARG REPO
COPY --from=templates /go/src/${REPO}/output/ /output/

# ---- addons output ----
FROM scratch AS addons-output
ARG REPO
COPY --from=addons /go/src/${REPO}/output/ /output/

# ---- charts output ----
FROM scratch AS patch-charts-output
ARG REPO
COPY --from=patch-charts /go/src/${REPO}/output/ /output/
