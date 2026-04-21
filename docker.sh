#!/bin/bash
# Makefile equivalent shell script for Docker operations

set -e

IMAGE_NAME="utbench"
CONTAINER_NAME="utbench-runner"

build() {
    echo "Building Docker image..."
    docker build -t ${IMAGE_NAME}:latest .
}

run() {
    echo "Running utbench in Docker..."
    docker run --rm \
        -v ./datasets:/app/datasets \
        -v ./artifacts:/app/artifacts \
        -v ./storage:/app/storage \
        -v ./configs:/app/configs \
        --env-file .env \
        ${IMAGE_NAME}:latest "$@"
}

run-python() {
    run run --models deepseek --langs python --max-samples 2 --dry-run
}

run-all() {
    run run --models deepseek --langs python,go,java,cpp --max-samples 5
}

shell() {
    echo "Opening shell in container..."
    docker run --rm -it \
        -v ./datasets:/app/datasets \
        -v ./artifacts:/app/artifacts \
        -v ./storage:/app/storage \
        --env-file .env \
        ${IMAGE_NAME}:latest /bin/bash
}

clean() {
    echo "Cleaning Docker resources..."
    docker rmi ${IMAGE_NAME}:latest 2>/dev/null || true
    docker rm ${CONTAINER_NAME} 2>/dev/null || true
}

case "$1" in
    build)
        build
        ;;
    run)
        run "${@:2}"
        ;;
    run-python)
        run-python
        ;;
    run-all)
        run-all
        ;;
    shell)
        shell
        ;;
    clean)
        clean
        ;;
    *)
        echo "Usage: $0 {build|run|run-python|run-all|shell|clean}"
        echo "  build        - Build Docker image"
        echo "  run <args>   - Run with custom arguments"
        echo "  run-python   - Quick test with Python only (dry-run)"
        echo "  run-all      - Run all languages"
        echo "  shell        - Open interactive shell in container"
        echo "  clean        - Remove Docker image and containers"
        exit 1
        ;;
esac