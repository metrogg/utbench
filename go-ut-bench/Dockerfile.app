# Fast application-layer rebuild.
# Requires a previously built toolchain image, normally utbench:latest.
#
# Usage:
#   docker build -f Dockerfile.app --build-arg BASE_IMAGE=utbench:latest -t utbench:latest .

ARG BASE_IMAGE=utbench:latest
FROM ${BASE_IMAGE}

WORKDIR /app

# Keep module download cached when only application code changes.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o utbench ./cmd/utbench/

RUN mkdir -p /app/artifacts /app/storage

ENTRYPOINT ["./utbench"]
CMD ["--help"]
