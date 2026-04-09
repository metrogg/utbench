#!/usr/bin/env bash

set -euo pipefail

PYTHON_BIN="${PYTHON_BIN:-python3}"

RESULTS_ROOT="results"
MODEL=""
LANG=""
MAX_SAMPLES=""
SAMPLE_GLOB=""
MAX_WORKERS=""
DRY_RUN="false"
NO_RESUME="false"
RESET_CHECKPOINT="false"
ALL_HISTORY="false"
SKIP_EVALUATOR="false"
SKIP_REPORTER="false"
ONLY_SELF_CONTAINED="false"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --results-root)
      RESULTS_ROOT="$2"
      shift 2
      ;;
    --model)
      MODEL="$2"
      shift 2
      ;;
    --lang)
      LANG="$2"
      shift 2
      ;;
    --max-samples)
      MAX_SAMPLES="$2"
      shift 2
      ;;
    --sample-glob)
      SAMPLE_GLOB="$2"
      shift 2
      ;;
    --max-workers)
      MAX_WORKERS="$2"
      shift 2
      ;;
    --dry-run)
      DRY_RUN="true"
      shift
      ;;
    --no-resume)
      NO_RESUME="true"
      shift
      ;;
    --reset-checkpoint)
      RESET_CHECKPOINT="true"
      shift
      ;;
    --all-history)
      ALL_HISTORY="true"
      shift
      ;;
    --skip-evaluator)
      SKIP_EVALUATOR="true"
      shift
      ;;
    --skip-reporter)
      SKIP_REPORTER="true"
      shift
      ;;
    --only-self-contained)
      ONLY_SELF_CONTAINED="true"
      shift
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

runner_cmd=("$PYTHON_BIN" -m benchmark.runner --results-root "$RESULTS_ROOT")

if [[ -n "$MODEL" ]]; then
  runner_cmd+=(--model "$MODEL")
fi
if [[ -n "$LANG" ]]; then
  runner_cmd+=(--lang "$LANG")
fi
if [[ -n "$MAX_SAMPLES" ]]; then
  runner_cmd+=(--max-samples "$MAX_SAMPLES")
fi
if [[ -n "$SAMPLE_GLOB" ]]; then
  runner_cmd+=(--sample-glob "$SAMPLE_GLOB")
fi
if [[ -n "$MAX_WORKERS" ]]; then
  runner_cmd+=(--max-workers "$MAX_WORKERS")
fi
if [[ "$DRY_RUN" == "true" ]]; then
  runner_cmd+=(--dry-run)
fi
if [[ "$NO_RESUME" == "true" ]]; then
  runner_cmd+=(--no-resume)
fi
if [[ "$RESET_CHECKPOINT" == "true" ]]; then
  runner_cmd+=(--reset-checkpoint)
fi

echo "[run_benchmark] runner: ${runner_cmd[*]}"
"${runner_cmd[@]}"

if [[ "$SKIP_EVALUATOR" == "true" ]]; then
  echo "[run_benchmark] evaluator skipped"
  exit 0
fi

evaluator_cmd=("$PYTHON_BIN" -m benchmark.evaluator --results-root "$RESULTS_ROOT")
if [[ -n "$MODEL" ]]; then
  evaluator_cmd+=(--model "$MODEL")
fi
if [[ -n "$LANG" ]]; then
  evaluator_cmd+=(--lang "$LANG")
fi
if [[ "$ALL_HISTORY" == "true" ]]; then
  evaluator_cmd+=(--all-history)
fi
if [[ "$ONLY_SELF_CONTAINED" == "true" ]]; then
  evaluator_cmd+=(--only-self-contained)
fi

echo "[run_benchmark] evaluator: ${evaluator_cmd[*]}"
"${evaluator_cmd[@]}"

if [[ "$SKIP_REPORTER" == "true" ]]; then
  echo "[run_benchmark] reporter skipped"
  exit 0
fi

reporter_cmd=("$PYTHON_BIN" -m benchmark.reporter --results-root "$RESULTS_ROOT")
echo "[run_benchmark] reporter: ${reporter_cmd[*]}"
"${reporter_cmd[@]}"
