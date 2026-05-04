# UT-Bench Agent Sandbox: Universal image (all agents, all languages)
# Contains: OpenCode CLI + CodeBuddy CLI + Claude Code CLI + Python + Go + Java + C++ + test frameworks.
# Does NOT contain: coverage / mutation tools (those live in the eval image).
#
# This is a unified sandbox image shared by all CLI Agent frameworks.
# New agents only need to add their CLI installation in Stage 1.

# ── Stage 1: Install all Agent CLIs ──────────────────────────
FROM node:22-bookworm-slim AS agent-clis

ARG NPM_REGISTRY=
RUN if [ -n "${NPM_REGISTRY}" ]; then npm config set registry "${NPM_REGISTRY}"; fi

# OpenCode CLI
RUN npm install -g opencode-ai opencode-linux-x64-baseline

# CodeBuddy CLI (腾讯)
RUN npm install -g @tencent-ai/codebuddy-code

# Claude Code CLI (Anthropic)
RUN npm install -g @anthropic-ai/claude-code

# ── Stage 2: Full runtime image ─────────────────────────────
FROM python:3.11-slim-bookworm

ARG DEBIAN_MIRROR=
ARG PIP_INDEX_URL=
ARG GO_VERSION=1.22.12
ARG GO_DOWNLOAD_URLS=

# Copy all Agent CLIs from stage 1
COPY --from=agent-clis /usr/local/bin/ /usr/local/bin/
COPY --from=agent-clis /usr/local/lib/ /usr/local/lib/

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

# ── 非 root 用户（Claude Code 拒绝以 root 使用 --dangerously-skip-permissions）──
RUN useradd -m -s /bin/bash agent \
    && mkdir -p /home/agent/.npm /home/agent/.cache \
    && chown -R agent:agent /home/agent

WORKDIR /workspace
