# go-ut-bench Docker Image
# Ubuntu 24.04 with all evaluation toolchains pre-installed
# Supports: Python, Go, Java, C++ unit test evaluation

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

# System dependencies - Ubuntu 24.04 provides LLVM 18 natively
RUN apt-get update && apt-get install -y --fix-missing --no-install-recommends \
    build-essential \
    gcc \
    g++ \
    cmake \
    git \
    curl \
    wget \
    unzip \
    python3 \
    python3-pip \
    python3-dev \
    openjdk-21-jdk \
    maven \
    libgtest-dev \
    gdb \
    sudo \
    gnupg \
    lsb-release \
    software-properties-common \
    clang-18 \
    llvm-18 \
    llvm-18-dev \
    libclang-18-dev \
    lld-18 \
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

# LLVM/Clang 18 paths
ENV LLVM_DIR=/usr/lib/llvm-18
ENV PATH="${LLVM_DIR}/bin:${PATH}"
ENV LD_LIBRARY_PATH="${LLVM_DIR}/lib:/usr/lib/x86_64-linux-gnu:/lib/x86_64-linux-gnu:/usr/lib64"
ENV LLVM_CONFIG_PATH=/usr/bin/llvm-config-18

# Install Mull 19 (C++ mutation testing) from Cloudsmith
RUN curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && \
    apt-get update && apt-get install -y mull-19 && \
    rm -rf /var/lib/apt/lists/*

# Install Python testing tools
RUN pip3 install pytest coverage mutmut

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