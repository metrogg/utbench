# go-ut-bench Docker Image
# Ubuntu 24.04 with all evaluation toolchains pre-installed
# Supports: Python, Go, Java, C++ unit test evaluation

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

<<<<<<< HEAD
# Use Aliyun mirror for faster apt in China
RUN sed -i 's|http://archive.ubuntu.com/ubuntu|http://mirrors.aliyun.com/ubuntu|g' /etc/apt/sources.list.d/ubuntu.sources && \
    sed -i 's|http://security.ubuntu.com/ubuntu|http://mirrors.aliyun.com/ubuntu|g' /etc/apt/sources.list.d/ubuntu.sources

# Install base tools first
=======
# Install base tools first (needed for adding LLVM repo)
>>>>>>> origin/feat/go
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    wget \
    gnupg \
    lsb-release \
    software-properties-common \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

<<<<<<< HEAD
# Add LLVM official apt repository for LLVM 19 on Ubuntu 24.04
=======
# Add LLVM official repository for LLVM 19
>>>>>>> origin/feat/go
RUN wget -qO- https://apt.llvm.org/llvm-snapshot.gpg.key | tee /etc/apt/trusted.gpg.d/apt.llvm.org.asc && \
    add-apt-repository "deb http://apt.llvm.org/noble/ llvm-toolchain-noble-19 main" && \
    apt-get update

<<<<<<< HEAD
# System dependencies
RUN apt-get install -y --no-install-recommends \
=======
# System dependencies - Install LLVM 19 for Mull 19 compatibility
RUN apt-get update && apt-get install -y --fix-missing --no-install-recommends \
>>>>>>> origin/feat/go
    build-essential \
    gcc \
    g++ \
    cmake \
    git \
    unzip \
    python3 \
    python3-pip \
    python3-dev \
<<<<<<< HEAD
=======
    python3-pytest \
    python3-six \
    python3-dateutil \
>>>>>>> origin/feat/go
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
<<<<<<< HEAD
    && cd /usr/src/googletest && cmake . && make && cp lib/*.a /usr/lib/ \
    && rm -rf /var/lib/apt/lists/*

# Install pytest via pip (more reliable than apt package naming)
RUN pip3 install --break-system-packages pytest

=======
    && cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/ \
    && rm -rf /var/lib/apt/lists/*

>>>>>>> origin/feat/go
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

<<<<<<< HEAD
# Install Mull 19 (C++ mutation testing) from Cloudsmith
RUN curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && \
    apt-get update && apt-get install -y mull-19 && \
=======
# 配置系统动态库搜索路径（解决 Mull "Could not find dynamic library" warning）
RUN echo "/usr/lib/llvm-19/lib" > /etc/ld.so.conf.d/llvm-19.conf && \
    echo "/usr/lib/x86_64-linux-gnu" >> /etc/ld.so.conf.d/llvm-19.conf && \
    echo "/lib/x86_64-linux-gnu" >> /etc/ld.so.conf.d/llvm-19.conf && \
    ldconfig

# Install Mull 19 (C++ mutation testing) from Cloudsmith
RUN curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && \
    apt-get update && apt-get install -y mull-19 && \
    ldconfig && \
>>>>>>> origin/feat/go
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
<<<<<<< HEAD
CMD ["--help"]
=======
CMD ["--help"]
>>>>>>> origin/feat/go
