# AGENTS.md

## 1) Purpose

Evaluate multi-model unit-test generation quality across languages.

Current state:
- **Python**: full pipeline (compile → test → coverage → mutation) works
- **Java/Go/Cpp/JavaScript**: evaluator has placeholders only, no real toolchain
- **Runner**: works for all languages; evaluator only processes Python
- **Mutation testing**: requires Linux (mutmut not validated on Windows)

## 2) Key Paths

- `benchmark/runner/` — model API calls, prompt construction, response parsing
- `benchmark/evaluator/` — quality evaluation (Python only functional)
- `benchmark/reporter/` — generates JSON/CSV/HTML reports from evaluator output
- `benchmark/config/models.yaml` — model and benchmark configuration
- `dataset/` — benchmark samples organized by `<lang>/<complexity>/`
- `results/` — generated tests, artifacts, checkpoints, summaries

## 3) CLIs

### Runner

```bash
python -m benchmark.runner                          # all enabled models × all languages
python -m benchmark.runner --model deepseek --lang python --max-samples 2
python -m benchmark.runner --model deepseek --lang python --sample-glob "boundary/*.py"
python -m benchmark.runner --dry-run                # no real API calls
python -m benchmark.runner --reset-checkpoint       # clear checkpoint before run
python -m benchmark.runner --no-resume              # force full rerun
```

### Evaluator

```bash
python -m benchmark.evaluator --results-root results
python -m benchmark.evaluator --results-root results --model deepseek --lang python
python -m benchmark.evaluator --results-root results --only-self-contained  # Python allowlist only
python -m benchmark.evaluator --results-root results --all-history          # all test files per sample
```

### Reporter

```bash
python -m benchmark.reporter --results-root results
python -m benchmark.reporter --results-root results --evaluator-summary results/evaluator_summary_*.json
```

### Pipeline script (runner → evaluator → reporter)

```bash
bash scripts/run_benchmark.sh
bash scripts/run_benchmark.sh --model deepseek --lang python --skip-reporter
```

## 4) Environment Setup

```bash
pip install pyyaml pytest coverage mutmut
bash scripts/setup.sh
```

API keys (set before running runner):

```powershell
$env:DASHSCOPE_API_KEY="..."
$env:DEEPSEEK_API_KEY="..."
$env:MINIMAX_API_KEY="..."
$env:VOLCENGINE_API_KEY="..."
$env:BIGMODEL_API_KEY="..."   # for disabled 'glm' model
```

## 5) Model Config (models.yaml)

Required fields per model: `enabled`, `provider`, `config.api_endpoint`, `config.model`, `config.api_key_env`

Enabled models: `qwen`, `deepseek`, `minimax`, `doubao-seed`, `glm-4.7`

## 6) Language Inference (Runner/Evaluator Contract)

Runner infers language from dataset path `dataset/<lang>/...`, **not** file extension.

Evaluator loader supports legacy bug where scenario folders were used as language:
- `boundary`, `simple_function`, `complex_dependency`, `interface_mock` → mapped to `python`

Generated test file pattern: `<model>_<language>_<sample_id>_<timestamp>.test.<ext>`

Source path resolution: metadata `sample_path` > dataset search by `sample_id`

## 7) Output Structure

```
results/
  <model>/
    tests/      # generated test files
    reports/    # metadata JSON
    artifacts/  # raw API responses
  checkpoints/  # resume state
  runner_summary_*.json
  evaluator_summary_*.json
  reports/reporter_*.{json,csv,html}
```

## 8) Code Conventions

- User-facing docs in Chinese; code comments in English when needed
- Use `pathlib.Path` for filesystem operations
- Keep `None` for unavailable metrics (never fake zeros in summary JSON)
- Follow PEP 8, type hints for public functions

## 9) Gotchas

- No unit tests exist for benchmark code itself (`tests/` is empty)
- Evaluator `--only-self-contained` uses `benchmark/config/python_self_contained_allowlist.txt`
- Windows: mutation testing skipped (mutmut needs Linux)
- `go-ut-bench/` and `datasets/` directories are separate projects, not main benchmark
