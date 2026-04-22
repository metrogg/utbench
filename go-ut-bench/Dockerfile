# go-ut-bench Docker Image
# Ubuntu 24.04 with all evaluation toolchains pre-installed
# Supports: Python, Go, Java, C++ unit test evaluation

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

# Install base tools first (needed for adding LLVM repo)
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    gnupg \
    lsb-release \
    software-properties-common \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Add LLVM official repository for LLVM 19
RUN wget -qO- https://apt.llvm.org/llvm-snapshot.gpg.key | tee /etc/apt/trusted.gpg.d/apt.llvm.org.asc && \
    add-apt-repository "deb http://apt.llvm.org/noble/ llvm-toolchain-noble-19 main" && \
    apt-get update

# System dependencies - Install LLVM 19 for Mull 19 compatibility
RUN apt-get update && apt-get install -y --fix-missing --no-install-recommends \
    build-essential \
    gcc \
    g++ \
    cmake \
    git \
    unzip \
    python3 \
    python3-pip \
    python3-dev \
    python3-pytest \
    python3-venv \
    openjdk-21-jdk \
    maven \
    libgtest-dev \
    gdb \
    sudo \
    clang-19 \
    llvm-19 \
    llvm-19-dev \
    libclang-19-dev \
    lld-19 \
    libstdc++6 \
    libgcc-s1 \
    && cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/ \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.24 (latest stable)
RUN wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz && \
    rm go1.24.2.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
ENV GOPROXY=https://goproxy.cn,direct

# LLVM/Clang 19 paths
ENV LLVM_DIR=/usr/lib/llvm-19
ENV PATH="${LLVM_DIR}/bin:${PATH}"
ENV LD_LIBRARY_PATH="${LLVM_DIR}/lib:/usr/lib/x86_64-linux-gnu:/lib/x86_64-linux-gnu:/usr/lib64"
ENV LLVM_CONFIG_PATH=/usr/bin/llvm-config-19

# Install Mull 19 (C++ mutation testing) from Cloudsmith
RUN curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && \
    apt-get update && apt-get install -y mull-19 && \
    rm -rf /var/lib/apt/lists/*

# Install Python testing tools - use venv for mutmut to avoid conflicts
RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install --upgrade pip && \
    /opt/venv/bin/pip install coverage mutmut && \
    ln -sf /opt/venv/bin/mutmut /usr/local/bin/mutmut && \
    ln -sf /opt/venv/bin/coverage /usr/local/bin/coverage

# Install Go mutation testing tool (gremlins - modern replacement for go-mutesting)
RUN go install github.com/go-gremlins/gremlins/cmd/gremlins@latest

WORKDIR /app

# Copy go.mod first for dependency cache (cached if go.mod unchanged)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build (this layer rebuilds when code changes)
COPY . .
RUN go build -o utbench ./cmd/utbench/

RUN mkdir -p /app/artifacts /app/storage

ENTRYPOINT ["./utbench"]
CMD ["--help"]