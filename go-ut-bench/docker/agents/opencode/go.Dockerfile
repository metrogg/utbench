FROM node:22-bookworm-slim AS opencode

ARG NPM_REGISTRY=
RUN if [ -n "${NPM_REGISTRY}" ]; then npm config set registry "${NPM_REGISTRY}"; fi \
    && npm install -g opencode-ai opencode-linux-x64-baseline

FROM golang:1.22-bookworm

ARG DEBIAN_MIRROR=

COPY --from=opencode /usr/local/bin/ /usr/local/bin/
COPY --from=opencode /usr/local/lib/ /usr/local/lib/

RUN if [ -n "${DEBIAN_MIRROR}" ]; then \
      if [ -f /etc/apt/sources.list.d/debian.sources ]; then \
        sed -i "s|http://deb.debian.org/debian|${DEBIAN_MIRROR}|g" /etc/apt/sources.list.d/debian.sources; \
      elif [ -f /etc/apt/sources.list ]; then \
        sed -i "s|http://deb.debian.org/debian|${DEBIAN_MIRROR}|g" /etc/apt/sources.list; \
      fi; \
    fi \
    && apt-get -o Acquire::Retries=5 -o Acquire::http::Timeout=60 update \
    && apt-get install -y --fix-missing --no-install-recommends \
        ca-certificates \
        git \
        ripgrep \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace
