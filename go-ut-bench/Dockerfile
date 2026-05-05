# go-ut-bench Docker Image
# Ubuntu 24.04 with all evaluation toolchains pre-installed
# Supports: Python, Go, Java, C++ unit test evaluation

FROM docker:28-cli AS docker-cli

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Shanghai

ARG UBUNTU_MIRROR=

# Optionally override Ubuntu archive/security hosts for build environments
# that need a regional mirror. Leave UBUNTU_MIRROR empty to keep upstream
# defaults and avoid hard failures when a third-party mirror is flaky.
# Note: only replace the host portion, keep http:// protocol to avoid
# certificate verification issues before ca-certificates is installed.
RUN if [ -n "${UBUNTU_MIRROR}" ]; then \
      MIRROR_HOST=$(echo "${UBUNTU_MIRROR}" | sed 's|https\?://||;s|/$||'); \
      if [ -f /etc/apt/sources.list.d/ubuntu.sources ]; then \
        sed -i "s|archive.ubuntu.com/ubuntu|${MIRROR_HOST}|g" /etc/apt/sources.list.d/ubuntu.sources && \
        sed -i "s|security.ubuntu.com/ubuntu|${MIRROR_HOST}|g" /etc/apt/sources.list.d/ubuntu.sources; \
      fi; \
      if [ -f /etc/apt/sources.list ]; then \
        sed -i "s|archive.ubuntu.com/ubuntu|${MIRROR_HOST}|g" /etc/apt/sources.list && \
        sed -i "s|security.ubuntu.com/ubuntu|${MIRROR_HOST}|g" /etc/apt/sources.list; \
      fi; \
    fi

# Install base tools first (needed for adding LLVM repo)
RUN apt-get -o Acquire::Retries=5 -o Acquire::http::Timeout=60 update && apt-get install -y --no-install-recommends \
    apt-transport-https \
    curl \
    wget \
    gnupg \
    lsb-release \
    software-properties-common \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# System dependencies. Keep apt-get update in the same layer as install so a
# cached repository setup layer never leaves this step with stale package lists.
RUN apt-get -o Acquire::Retries=5 -o Acquire::http::Timeout=60 update && \
    apt-get install -y --fix-missing --no-install-recommends \
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
    python3-six \
    python3-dateutil \
    python3-venv \
    default-jdk \
    maven \
    libgtest-dev \
    gdb \
    sudo \
    libstdc++6 \
    libgcc-s1 && \
    cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/

# Add LLVM official repository for LLVM 19 only after base dependencies are
# installed. This keeps temporary apt.llvm.org failures from breaking the large
# Ubuntu package layer above.
RUN wget -qO- https://apt.llvm.org/llvm-snapshot.gpg.key | tee /etc/apt/trusted.gpg.d/apt.llvm.org.asc && \
    add-apt-repository "deb http://apt.llvm.org/noble/ llvm-toolchain-noble-19 main"

# LLVM/Clang 19 is installed separately so failures from apt.llvm.org are
# easier to diagnose than a single very long apt-get install command.
RUN apt-get -o Acquire::Retries=5 -o Acquire::http::Timeout=60 update && \
    apt-get install -y --fix-missing --no-install-recommends \
    clang-19 \
    llvm-19 \
    llvm-19-dev \
    libclang-19-dev \
    lld-19 \
    && rm -rf /var/lib/apt/lists/*

# Configure Maven to use Aliyun mirror for faster downloads in China
RUN mkdir -p /root/.m2 && cat > /root/.m2/settings.xml << 'SETTINGSEOF'
<settings>
  <mirrors>
    <mirror>
      <id>aliyun</id>
      <mirrorOf>central</mirrorOf>
      <name>Aliyun Maven Mirror</name>
      <url>https://maven.aliyun.com/repository/central</url>
    </mirror>
  </mirrors>
</settings>
SETTINGSEOF

# Pre-download Maven dependencies used by Java evaluation (JUnit 5, JaCoCo, Surefire, PITest).
# This avoids downloading them on every container run, which is slow and unreliable
# with concurrent workers.
RUN mkdir -p /tmp/maven-warmup && cat > /tmp/maven-warmup/pom.xml << 'POMEOF'
<project>
  <modelVersion>4.0.0</modelVersion>
  <groupId>warmup</groupId><artifactId>warmup</artifactId><version>1.0</version>
  <properties>
    <maven.compiler.source>17</maven.compiler.source>
    <maven.compiler.target>17</maven.compiler.target>
  </properties>
  <dependencies>
    <dependency><groupId>org.junit.jupiter</groupId><artifactId>junit-jupiter</artifactId><version>5.10.2</version><scope>test</scope></dependency>
    <dependency><groupId>org.junit.vintage</groupId><artifactId>junit-vintage-engine</artifactId><version>5.10.2</version><scope>test</scope></dependency>
  </dependencies>
  <build><plugins>
    <plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-compiler-plugin</artifactId><version>3.11.0</version></plugin>
    <plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-surefire-plugin</artifactId><version>3.1.2</version></plugin>
    <plugin><groupId>org.jacoco</groupId><artifactId>jacoco-maven-plugin</artifactId><version>0.8.11</version></plugin>
    <plugin><groupId>org.pitest</groupId><artifactId>pitest-maven</artifactId><version>1.19.6</version>
      <dependencies><dependency><groupId>org.pitest</groupId><artifactId>pitest-junit5-plugin</artifactId><version>1.2.1</version></dependency></dependencies>
    </plugin>
  </plugins></build>
</project>
POMEOF
RUN cd /tmp/maven-warmup && mvn dependency:resolve test-compile -q && rm -rf /tmp/maven-warmup

# Copy docker CLI binary from official image (lighter than docker.io package)
COPY --from=docker-cli /usr/local/bin/docker /usr/local/bin/docker

# Install Go 1.24. Keep multiple download endpoints because TLS to go.dev can
# be flaky in some build networks.
ARG GO_VERSION=1.24.2
ARG GO_DOWNLOAD_URLS=
RUN set -eux; \
    tarball="go${GO_VERSION}.linux-amd64.tar.gz"; \
    urls="${GO_DOWNLOAD_URLS:-https://go.dev/dl/${tarball} https://dl.google.com/go/${tarball} https://mirrors.aliyun.com/golang/${tarball}}"; \
    for url in ${urls}; do \
        echo "Downloading Go from ${url}"; \
        if curl -fL --retry 5 --retry-delay 3 --retry-all-errors --connect-timeout 20 -o /tmp/go.tar.gz "${url}"; then \
            break; \
        fi; \
        rm -f /tmp/go.tar.gz; \
    done; \
    test -s /tmp/go.tar.gz; \
    rm -rf /usr/local/go; \
    tar -C /usr/local -xzf /tmp/go.tar.gz; \
    rm /tmp/go.tar.gz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
ENV GOPROXY=https://goproxy.cn,direct

# LLVM/Clang 19 paths
ENV LLVM_DIR=/usr/lib/llvm-19
ENV PATH="${LLVM_DIR}/bin:${PATH}"
ENV LD_LIBRARY_PATH="${LLVM_DIR}/lib:/usr/lib/x86_64-linux-gnu:/lib/x86_64-linux-gnu:/usr/lib64"
ENV LLVM_CONFIG_PATH=/usr/bin/llvm-config-19

# 配置系统动态库搜索路径（解决 Mull "Could not find dynamic library" warning）
RUN echo "/usr/lib/llvm-19/lib" > /etc/ld.so.conf.d/llvm-19.conf && \
    echo "/usr/lib/x86_64-linux-gnu" >> /etc/ld.so.conf.d/llvm-19.conf && \
    echo "/lib/x86_64-linux-gnu" >> /etc/ld.so.conf.d/llvm-19.conf && \
    ldconfig

# Install Mull (C++ mutation testing) from GitHub releases
# Package naming: Mull-{LLVM_VER}-{MULL_VER}-LLVM-{FULL_LLVM}-ubuntu-amd64-{UBUNTU_VER}.deb
ARG MULL_VERSION=0.33.0
ARG LLVM_MAJOR=19
RUN set -eux; \
    ARCH=$(dpkg --print-architecture); \
    UBUNTU_VER=$(lsb_release -rs); \
    BASE_URL="https://github.com/mull-project/mull/releases/download/${MULL_VERSION}"; \
    for ubu in "${UBUNTU_VER}" "24.04" "26.04"; do \
      for llvm_full in "19.1.1" "19.1.7" "19.1.6" "19.1.0"; do \
        name="Mull-${LLVM_MAJOR}-${MULL_VERSION}-LLVM-${llvm_full}-ubuntu-${ARCH}-${ubu}.deb"; \
        echo "Trying: ${BASE_URL}/${name}"; \
        if curl -fL --retry 3 --connect-timeout 10 "${BASE_URL}/${name}" -o /tmp/mull.deb 2>/dev/null; then \
          dpkg -i /tmp/mull.deb && rm /tmp/mull.deb && ldconfig; \
          exit 0; \
        fi; \
      done; \
    done; \
    echo "ERROR: Could not download Mull ${MULL_VERSION} for LLVM ${LLVM_MAJOR} on Ubuntu ${UBUNTU_VER}"; \
    exit 1

# Install Python testing tools - use venv for mutmut to avoid conflicts
RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install --upgrade pip && \
    /opt/venv/bin/pip install coverage mutmut && \
    ln -sf /opt/venv/bin/mutmut /usr/local/bin/mutmut && \
    ln -sf /opt/venv/bin/coverage /usr/local/bin/coverage

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
