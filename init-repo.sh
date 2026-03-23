#!/bin/bash
# init-repo.sh
# 用于将 ut-bench 目录结构初始化提交到 dev/init-ut-bench 分支
# 使用方式：
#   1. 将本脚本放到与 ut-bench 目录同级位置，或在 ut-bench 目录内执行
#   2. 确保已配置好 git remote（git remote add origin https://git.tencent.com/ut/ut-bench.git）
#   3. bash init-repo.sh

set -e

BRANCH="dev/init-ut-bench"
REMOTE="origin"

echo "=== ut-bench 仓库初始化 ==="

# 初始化 git（如果还没有）
if [ ! -d ".git" ]; then
    git init
    git remote add origin https://git.tencent.com/ut/ut-bench.git
    echo "[INFO] git 仓库已初始化"
fi

# 配置用户信息（可选）
# git config user.name "andrewjiang"
# git config user.email "andrewjiang@tencent.com"

# 拉取远程（如果远程已有内容）
git fetch origin 2>/dev/null || echo "[INFO] 远程仓库为空或无法连接，跳过 fetch"

# 创建并切换到目标分支
git checkout -b "$BRANCH" 2>/dev/null || git checkout "$BRANCH"

# 添加所有文件
git add .

# 提交
git commit -m "feat: init ut-bench repo structure and README

- 初始化代码库目录结构
- 添加项目 README（评测目标、指标体系、参评模型、快速开始）
- 添加数据集说明（dataset/README.md）
- 添加评测模块说明（benchmark/README.md）
- 添加模型配置文件（benchmark/config/models.yaml）
- 添加方案设计文档（docs/design/）
  - benchmark-design.md：评测方案总体设计
  - metrics-definition.md：评测指标定义
- 添加工具脚本（scripts/）
- 添加 .gitignore"

# 推送到远程
git push -u "$REMOTE" "$BRANCH"

echo ""
echo "=== 提交完成 ==="
echo "分支：$BRANCH"
echo "仓库：https://git.tencent.com/ut/ut-bench"
