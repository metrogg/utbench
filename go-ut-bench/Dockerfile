# go-ut-bench Docker Image
# Ubuntu 22.04 with all evaluation toolchains pre-installed
# Supports: Python, Go, Java, C++ unit test evaluation

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

# System dependencies (cached layer)
RUN apt-get update && apt-get install -y \
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
    openjdk-17-jdk \
    maven \
    libgtest-dev \
    gdb \
    sudo \
    gnupg \
    lsb-release \
    software-properties-common \
    clang-15 \
    llvm-15 \
    lld-15 \
    && cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/ \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.26 (required for go-mutesting)
RUN wget https://go.dev/dl/go1.26.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.26.2.linux-amd64.tar.gz && \
    rm go1.26.2.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
ENV GOPROXY=https://goproxy.cn,direct

# LLVM/Clang 15 paths (installed in system dependencies layer)
ENV LLVM_DIR=/usr/lib/llvm-15
ENV PATH="${LLVM_DIR}/bin:${PATH}"

# Install Mull 15 (C++ mutation testing) from Cloudsmith
RUN curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && \
    apt-get update && apt-get install -y mull-15 && \
    rm -rf /var/lib/apt/lists/*

# Install Python testing tools (cached layer)
RUN pip3 install pytest coverage mutmut

# Install Go mutation testing tool
RUN go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest

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