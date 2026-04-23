# Windows PowerShell 脚本
# 使用方法: .\run_bench.ps1

param(
    [string]$Models = "deepseek,minimax",
    [string]$Langs = "python,go,java,cpp",
    [int]$MaxSamples = 2,
    [bool]$Mutation = $true
)

# 检查密钥文件
if (-not (Test-Path ".env")) {
    Write-Host "错误: .env 文件不存在"
    exit 1
}

Write-Host "=========================================="
Write-Host "开始评测: $Models"
Write-Host "语言: $Langs"
Write-Host "每场景样本数: $MaxSamples"
Write-Host "变异测试: $Mutation"
Write-Host "=========================================="

docker run --rm `
    --env-file .env `
    -v "${PWD}/datasets:/app/datasets" `
    -v "${PWD}/artifacts:/app/artifacts" `
    -v "${PWD}/configs:/app/configs" `
    utbench:latest run `
        --models $Models `
        --langs $Langs `
        --dataset-root /app/datasets `
        --output-root /app/artifacts `
        --config /app/configs/models.yaml `
        --max-samples $MaxSamples `
        --mutation-enabled=$Mutation

Write-Host "=========================================="
Write-Host "完成! 结果在 artifacts/ 目录"
Write-Host "=========================================="