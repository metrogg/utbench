# ============================================================
# UT-Bench 统一构建脚本 (Windows PowerShell)
# 用法:
#   .\build.ps1 all              # 构建全部镜像（评测 + 沙箱）
#   .\build.ps1 eval             # 仅构建评测镜像
#   .\build.ps1 sandbox          # 构建统一 Agent 沙箱镜像
#   .\build.ps1 verify           # 验证所有镜像是否就绪
#
# 国内网络加速:
#   .\build.ps1 all -Cn          # 使用国内镜像源
# ============================================================

param(
    [Parameter(Position=0)]
    [string]$Target = "help",
    [switch]$Cn
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

# --- 国内镜像参数 ---
$EvalArgs = @()
$AgentArgs = @()

if ($Cn) {
    $EvalArgs = @(
        "--build-arg", "UBUNTU_MIRROR=https://mirrors.tuna.tsinghua.edu.cn/ubuntu",
        "--build-arg", "GO_DOWNLOAD_URLS=https://mirrors.aliyun.com/golang/go1.24.2.linux-amd64.tar.gz"
    )
    $AgentArgs = @(
        "--build-arg", "NPM_REGISTRY=https://registry.npmmirror.com",
        "--build-arg", "DEBIAN_MIRROR=https://mirrors.tuna.tsinghua.edu.cn/debian",
        "--build-arg", "PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple",
        "--build-arg", "GO_DOWNLOAD_URLS=https://mirrors.aliyun.com/golang/go1.22.12.linux-amd64.tar.gz"
    )
}

# --- 构建函数 ---
function Build-Eval {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "构建评测镜像: utbench:latest"
    Write-Host "==========================================" -ForegroundColor Cyan
    docker build @EvalArgs -t utbench:latest .
    if ($LASTEXITCODE -ne 0) { throw "评测镜像构建失败" }
    Write-Host "✅ utbench:latest 构建完成" -ForegroundColor Green
    Write-Host ""
}

function Build-Sandbox {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "构建统一 Agent 沙箱: utbench-agent-base:latest"
    Write-Host "==========================================" -ForegroundColor Cyan
    docker build @AgentArgs -f docker/agents/Dockerfile -t utbench-agent-base:latest docker/agents/
    if ($LASTEXITCODE -ne 0) { throw "沙箱镜像构建失败" }
    Write-Host "✅ utbench-agent-base:latest 构建完成" -ForegroundColor Green
    Write-Host ""
}

function Verify-Images {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "验证镜像状态"
    Write-Host "==========================================" -ForegroundColor Cyan
    $AllOk = $true

    $Images = @(
        "utbench:latest",
        "utbench-agent-base:latest"
    )

    foreach ($img in $Images) {
        $inspect = docker image inspect $img 2>$null
        if ($LASTEXITCODE -eq 0) {
            $size = (docker image inspect $img --format '{{.Size}}' 2>$null)
            $sizeMB = [math]::Round([int64]$size / 1MB)
            Write-Host "  ✅ $img (${sizeMB} MB)" -ForegroundColor Green
        } else {
            Write-Host "  ❌ $img (未构建)" -ForegroundColor Red
            $AllOk = $false
        }
    }
    Write-Host ""

    if ($AllOk) {
        Write-Host "✅ 全部镜像就绪" -ForegroundColor Green
    } else {
        Write-Host "⚠️  部分镜像缺失，运行 .\build.ps1 all 构建" -ForegroundColor Yellow
    }
}

# --- 主逻辑 ---
switch ($Target) {
    "all" {
        Build-Eval
        Build-Sandbox
        Verify-Images
    }
    "eval" { Build-Eval }
    "sandbox" { Build-Sandbox }
    "verify" { Verify-Images }
    default {
        Write-Host "UT-Bench 构建脚本"
        Write-Host ""
        Write-Host "用法: .\build.ps1 <target> [-Cn]"
        Write-Host ""
        Write-Host "Target:"
        Write-Host "  all         构建全部镜像（评测 + 沙箱）"
        Write-Host "  eval        仅构建评测镜像 (utbench:latest)"
        Write-Host "  sandbox     仅构建统一 Agent 沙箱 (utbench-agent-base:latest)"
        Write-Host "  verify      验证镜像状态"
        Write-Host ""
        Write-Host "选项:"
        Write-Host "  -Cn         使用国内镜像源加速构建"
        Write-Host ""
        Write-Host "示例:"
        Write-Host "  .\build.ps1 all -Cn     # 国内网络构建全部"
        Write-Host "  .\build.ps1 sandbox      # 只构建沙箱镜像"
        Write-Host "  .\build.ps1 verify       # 查看镜像状态"
    }
}
