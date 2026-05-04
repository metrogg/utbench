#!/bin/bash
set -euo pipefail

# ============================================================
# UT-Bench 统一构建脚本
# 用法:
#   ./build.sh all              # 构建全部镜像（评测 + 沙箱）
#   ./build.sh eval             # 仅构建评测镜像
#   ./build.sh sandbox          # 构建统一 Agent 沙箱镜像
#   ./build.sh verify           # 验证所有镜像是否就绪
#
# 国内网络加速:
#   ./build.sh all --cn         # 使用国内镜像源
# ============================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# --- 默认值 ---
USE_CN_MIRROR=false

# --- 解析参数 ---
TARGET="${1:-help}"
shift || true
for arg in "$@"; do
    case "$arg" in
        --cn) USE_CN_MIRROR=true ;;
        *) echo "未知参数: $arg"; exit 1 ;;
    esac
done

# --- 国内镜像参数 ---
BUILD_ARGS_EVAL=""
BUILD_ARGS_AGENT=""

if $USE_CN_MIRROR; then
    BUILD_ARGS_EVAL="--build-arg UBUNTU_MIRROR=https://mirrors.tuna.tsinghua.edu.cn/ubuntu --build-arg GO_DOWNLOAD_URLS=https://mirrors.aliyun.com/golang/go1.24.2.linux-amd64.tar.gz"
    BUILD_ARGS_AGENT="--build-arg NPM_REGISTRY=https://registry.npmmirror.com --build-arg DEBIAN_MIRROR=https://mirrors.tuna.tsinghua.edu.cn/debian --build-arg PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple --build-arg GO_DOWNLOAD_URLS=https://mirrors.aliyun.com/golang/go1.22.12.linux-amd64.tar.gz"
fi

# --- 构建函数 ---
build_eval() {
    echo "=========================================="
    echo "构建评测镜像: utbench:latest"
    echo "=========================================="
    docker build $BUILD_ARGS_EVAL -t utbench:latest .
    echo "✅ utbench:latest 构建完成"
    echo ""
}

build_sandbox() {
    echo "=========================================="
    echo "构建统一 Agent 沙箱: utbench-agent-opencode:latest"
    echo "=========================================="
    docker build $BUILD_ARGS_AGENT -f docker/agents/opencode/Dockerfile -t utbench-agent-opencode:latest docker/agents/opencode/
    echo "✅ utbench-agent-opencode:latest 构建完成"
    echo ""
}

verify_images() {
    echo "=========================================="
    echo "验证镜像状态"
    echo "=========================================="
    local all_ok=true

    check_image() {
        local name=$1
        if docker image inspect "$name" >/dev/null 2>&1; then
            local size
            size=$(docker image inspect "$name" --format '{{.Size}}' | awk '{printf "%.0f MB", $1/1024/1024}')
            echo "  ✅ $name ($size)"
        else
            echo "  ❌ $name (未构建)"
            all_ok=false
        fi
    }

    check_image "utbench:latest"
    check_image "utbench-agent-opencode:latest"
    echo ""

    if $all_ok; then
        echo "✅ 全部镜像就绪"
    else
        echo "⚠️  部分镜像缺失，运行 ./build.sh all 构建"
    fi
}

# --- 主逻辑 ---
case "$TARGET" in
    all)
        build_eval
        build_sandbox
        verify_images
        ;;
    eval)
        build_eval
        ;;
    sandbox)
        build_sandbox
        ;;
    verify)
        verify_images
        ;;
    help|*)
        echo "UT-Bench 构建脚本"
        echo ""
        echo "用法: ./build.sh <target> [--cn]"
        echo ""
        echo "Target:"
        echo "  all         构建全部镜像（评测 + 沙箱）"
        echo "  eval        仅构建评测镜像 (utbench:latest)"
        echo "  sandbox     仅构建统一 Agent 沙箱 (utbench-agent-opencode:latest)"
        echo "  verify      验证镜像状态"
        echo ""
        echo "选项:"
        echo "  --cn        使用国内镜像源加速构建"
        echo ""
        echo "示例:"
        echo "  ./build.sh all --cn     # 国内网络构建全部"
        echo "  ./build.sh sandbox      # 只构建沙箱镜像"
        echo "  ./build.sh verify       # 查看镜像状态"
        ;;
esac
