# ut-bench benchmark module

This directory contains benchmark-related modules:

- `runner/`: send dataset samples to model APIs and save generated tests
- `evaluator/`: quality evaluation pipeline
- `reporter/`: report generation module
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

## Evaluator quick start

Evaluate existing generated tests:

```bash
python -m benchmark.evaluator --results-root results
```

Scoped evaluation:

```bash
python -m benchmark.evaluator --results-root results --model deepseek --lang python
```

## Reporter quick start

Generate report artifacts (JSON/CSV/HTML) from evaluator outputs:

```bash
python -m benchmark.reporter --results-root results
```

Default HTML report is Chinese and visualized (summary + analysis + appendix), including:

- model comparison bar chart
- complexity trend line chart
- multi-dimension radar chart
- model-language heatmap

Specify evaluator summary file:

```bash
python -m benchmark.reporter --results-root results --evaluator-summary results/evaluator_summary_20260406T064045963256Z.json
```

Advanced options:

```bash
python -m benchmark.reporter \
  --results-root results \
  --formats json,csv,html \
  --chart-style teal \
  --threshold-test 0.7 \
  --threshold-line 0.7 \
  --threshold-branch 0.6 \
  --threshold-mutation 0.85
```

Shell wrapper:

```bash
bash scripts/gen_report.sh --results-root results
```

Reporter outputs follow:

```text
results/reports/
  reporter_summary_*.json
  reporter_by_model_*.csv
  reporter_by_language_*.csv
  reporter_samples_*.csv
  reporter_report_*.html
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
