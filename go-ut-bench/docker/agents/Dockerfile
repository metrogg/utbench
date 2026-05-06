# UT-Bench Agent Sandbox
#
# Purpose:
# - Run CLI agents such as Codex, Claude Code, OpenCode, and CodeBuddy.
# - Keep this image focused on generation-time tools so "Install Agent CLI"
#   is fast and is not blocked by evaluator toolchain packages.
#
# Evaluation compilers, coverage tools, and mutation tools live in the main
# evaluator image built from the repository root Dockerfile.

FROM node:22-bookworm-slim AS agent-clis

ARG NPM_REGISTRY=
ARG EXTRA_NPM_PACKAGES=

RUN if [ -n "${NPM_REGISTRY}" ]; then npm config set registry "${NPM_REGISTRY}"; fi

# Built-in runtimes used by the default agents.yaml.
RUN npm install -g opencode-ai opencode-linux-x64-baseline
RUN npm install -g @tencent-ai/codebuddy-code
RUN npm install -g @anthropic-ai/claude-code
# Pin to 0.42.x: 0.43+ removed `wire_api="chat"` support and forces /v1/responses,
# which国产兼容端点 (MiMo / DeepSeek / 豆包等) 暂不支持。
# 见 https://github.com/openai/codex/discussions/7782
RUN npm install -g @openai/codex@^0.42.0

# Optional Agent CLI selected from the Web installer, for example:
#   --build-arg EXTRA_NPM_PACKAGES="@openai/codex"
RUN if [ -n "${EXTRA_NPM_PACKAGES}" ]; then npm install -g ${EXTRA_NPM_PACKAGES}; fi

FROM node:22-bookworm-slim

ARG DEBIAN_MIRROR=

COPY --from=agent-clis /usr/local/bin/ /usr/local/bin/
COPY --from=agent-clis /usr/local/lib/ /usr/local/lib/

RUN set -eux; \
    if [ -n "${DEBIAN_MIRROR}" ]; then \
      MIRROR_HOST=$(echo "${DEBIAN_MIRROR}" | sed 's|https\?://||;s|/$||'); \
      if [ -f /etc/apt/sources.list.d/debian.sources ]; then \
        sed -i "s|deb.debian.org/debian|${MIRROR_HOST}|g" /etc/apt/sources.list.d/debian.sources; \
        sed -i "s|security.debian.org/debian-security|${MIRROR_HOST}-security|g" /etc/apt/sources.list.d/debian.sources || true; \
      elif [ -f /etc/apt/sources.list ]; then \
        sed -i "s|deb.debian.org/debian|${MIRROR_HOST}|g" /etc/apt/sources.list; \
        sed -i "s|security.debian.org/debian-security|${MIRROR_HOST}-security|g" /etc/apt/sources.list || true; \
      fi; \
    fi; \
    if apt-get -o Acquire::Retries=5 -o Acquire::http::Timeout=60 update; then \
      apt-get install -y --no-install-recommends \
        bash \
        ca-certificates \
        git \
        ripgrep \
        curl \
        wget \
        unzip \
        python3 \
        python3-pip \
      || echo "warning: optional agent helper tools failed to install; CLI runtime remains usable"; \
    else \
      echo "warning: apt update failed; skipping optional agent helper tools"; \
    fi; \
    rm -rf /var/lib/apt/lists/*

RUN useradd -m -s /bin/bash agent \
    && mkdir -p /home/agent/.npm /home/agent/.cache /workspace \
    && chown -R agent:agent /home/agent /workspace

WORKDIR /workspace
ENV PATH="/usr/local/bin:${PATH}"
