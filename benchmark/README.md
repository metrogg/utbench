# ut-bench benchmark module

This directory contains benchmark-related modules:

- `runner/`: send dataset samples to model APIs and save generated tests
- `evaluator/`: quality evaluation pipeline (to be implemented)
- `reporter/`: report generation module (to be implemented)
- `config/`: benchmark configuration (`models.yaml`)

## Runner quick start

Run a dry-run (no real API call):

```bash
python -m benchmark.runner --model deepseek --lang python --max-samples 2 --dry-run
```

Run with real model API:

```bash
python -m benchmark.runner --model deepseek --lang python --max-samples 2
```

Run with concurrency and checkpoint resume:

```bash
# default: resume enabled, use concurrency from benchmark.parallel in models.yaml
python -m benchmark.runner --model doubao-seed,kimi-k2.5,glm-4.7

# override total worker threads
python -m benchmark.runner --model doubao-seed,kimi-k2.5,glm-4.7 --max-workers 12

# disable resume and force full rerun
python -m benchmark.runner --model doubao-seed,kimi-k2.5,glm-4.7 --no-resume

# clear checkpoint for current scope then run
python -m benchmark.runner --model doubao-seed,kimi-k2.5,glm-4.7 --reset-checkpoint
```

Before real API calls, export the corresponding key in environment variables:

- `DEEPSEEK_API_KEY`
- `BIGMODEL_API_KEY`
- `DASHSCOPE_API_KEY`
- `MINIMAX_API_KEY`
- `VOLCENGINE_API_KEY`

Outputs follow:

```text
results/
  <model>/
    tests/      # generated test files
    reports/    # metadata json files
    artifacts/  # raw API responses and failure logs
```

## Prompt module design

The prompt module is in `benchmark/runner/prompt_builder.py`.

Responsibilities:

- Read test samples from `dataset/`
- Construct prompts with source code + test requirements + language markers
- Provide prompt text to Runner for model API calls
- Support stable response parsing by enforcing output format constraints

Prompt required components:

- Source code under test: loaded from dataset sample files
- Test requirement instructions: test framework, determinism, coverage targets
- Language identification: Java/Python/Go/cpp/javascript with framework mapping
- Context information: sample ID, scenario/complexity, dependency hints, mock guidance
