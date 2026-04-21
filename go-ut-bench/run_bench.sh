#!/bin/bash
# WSL / Linux 脚本
# 使用方法: ./run_bench.sh

MODELS="${1:-deepseek,minimax}"
LANGS="${2:-python,go,java,cpp}"
MAX_SAMPLES="${3:-2}"
MUTATION="${4:-true}"

# 检查密钥文件
if [ ! -f ".env" ]; then
    echo "错误: .env 文件不存在"
    exit 1
fi

mkdir -p artifacts

echo "=========================================="
echo "开始评测: $MODELS"
echo "语言: $LANGS"
echo "每场景样本数: $MAX_SAMPLES"
echo "变异测试: $MUTATION"
echo "=========================================="

docker run --rm \
    --env-file .env \
    -v "$(pwd)/datasets:/app/datasets" \
    -v "$(pwd)/artifacts:/app/artifacts" \
    -v "$(pwd)/configs:/app/configs" \
    utbench:latest run \
        --models "$MODELS" \
        --langs "$LANGS" \
        --dataset-root /app/datasets \
        --output-root /app/artifacts \
        --config /app/configs/models.yaml \
        --max-samples "$MAX_SAMPLES" \
        --mutation-enabled="$MUTATION"

echo "=========================================="
echo "完成! 结果在 artifacts/ 目录"
echo "=========================================="