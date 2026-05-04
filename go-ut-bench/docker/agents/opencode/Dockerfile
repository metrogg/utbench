# UT-Bench Agent Sandbox: Unified image (all languages)
# Used by CLI Agent subjects (OpenCode) to generate unit tests.
# Contains: OpenCode CLI + Python + Go + Java + C++ + test frameworks.
# Does NOT contain: coverage / mutation tools (those live in the eval image).

FROM node:22-bookworm-slim AS opencode

ARG NPM_REGISTRY=
RUN if [ -n "${NPM_REGISTRY}" ]; then npm config set registry "${NPM_REGISTRY}"; fi \
    && npm install -g opencode-ai opencode-linux-x64-baseline

# ── Base: Python + system tools ──────────────────────────────
FROM python:3.11-slim-bookworm

ARG DEBIAN_MIRROR=
ARG PIP_INDEX_URL=
ARG GO_VERSION=1.22.12
ARG GO_DOWNLOAD_URLS=

COPY --from=opencode /usr/local/bin/ /usr/local/bin/
COPY --from=opencode /usr/local/lib/ /usr/local/lib/

RUN if [ -n "${DEBIAN_MIRROR}" ]; then \
      MIRROR_HOST=$(echo "${DEBIAN_MIRROR}" | sed 's|https\?://||;s|/$||'); \
      if [ -f /etc/apt/sources.list.d/debian.sources ]; then \
        sed -i "s|deb.debian.org/debian|${MIRROR_HOST}|g" /etc/apt/sources.list.d/debian.sources; \
      elif [ -f /etc/apt/sources.list ]; then \
        sed -i "s|deb.debian.org/debian|${MIRROR_HOST}|g" /etc/apt/sources.list; \
      fi; \
    fi \
    && apt-get -o Acquire::Retries=5 -o Acquire::http::Timeout=60 update \
    && apt-get install -y --fix-missing --no-install-recommends \
        bash \
        ca-certificates \
        git \
        ripgrep \
        curl \
        wget \
        unzip \
        build-essential \
        gcc \
        g++ \
        cmake \
        # Java
        default-jdk-headless \
        maven \
        # C++ test framework
        libgtest-dev \
    # Python test framework
    && if [ -n "${PIP_INDEX_URL}" ]; then \
         pip install --no-cache-dir -i "${PIP_INDEX_URL}" pytest; \
       else \
         pip install --no-cache-dir pytest; \
       fi \
    # Build GoogleTest static library
    && cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/ \
    && rm -rf /var/lib/apt/lists/*

# ── Go ───────────────────────────────────────────────────────
RUN set -eux; \
    tarball="go${GO_VERSION}.linux-amd64.tar.gz"; \
    urls="${GO_DOWNLOAD_URLS:-https://go.dev/dl/${tarball} https://mirrors.aliyun.com/golang/${tarball}}"; \
    downloaded=false; \
    for url in $urls; do \
      curl -fL --retry 3 --retry-delay 2 --retry-all-errors -o /tmp/go.tar.gz "${url}" && downloaded=true && break || true; \
    done; \
    if ! $downloaded; then echo "Go download failed from all URLs" && exit 1; fi; \
    tar -C /usr/local -xzf /tmp/go.tar.gz; \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace
