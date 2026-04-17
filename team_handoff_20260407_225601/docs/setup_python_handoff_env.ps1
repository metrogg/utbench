$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent $PSScriptRoot
$VenvPython = Join-Path $ProjectRoot "venv\\Scripts\\python.exe"
$Requirements = Join-Path $ProjectRoot "env\\requirements_python_project_eval.txt"

if (-not (Test-Path -LiteralPath (Join-Path $ProjectRoot "venv"))) {
    python -m venv (Join-Path $ProjectRoot "venv")
}

& $VenvPython -m pip install --upgrade pip
& $VenvPython -m pip install -r $Requirements

Write-Host "Python handoff environment is ready."
