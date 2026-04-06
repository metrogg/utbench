#!/usr/bin/env bash

set -euo pipefail

RESULTS_ROOT="results"
EVALUATOR_SUMMARY=""
OUTPUT_DIR=""
FORMATS="json,csv,html"
CHART_STYLE="teal"
TOP_N_ERRORS="20"
THRESHOLD_COMPILE="1.0"
THRESHOLD_TEST="0.7"
THRESHOLD_LINE="0.7"
THRESHOLD_BRANCH="0.6"
THRESHOLD_MUTATION="0.85"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --results-root)
      RESULTS_ROOT="$2"
      shift 2
      ;;
    --evaluator-summary)
      EVALUATOR_SUMMARY="$2"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --formats)
      FORMATS="$2"
      shift 2
      ;;
    --chart-style)
      CHART_STYLE="$2"
      shift 2
      ;;
    --top-n-errors)
      TOP_N_ERRORS="$2"
      shift 2
      ;;
    --threshold-compile)
      THRESHOLD_COMPILE="$2"
      shift 2
      ;;
    --threshold-test)
      THRESHOLD_TEST="$2"
      shift 2
      ;;
    --threshold-line)
      THRESHOLD_LINE="$2"
      shift 2
      ;;
    --threshold-branch)
      THRESHOLD_BRANCH="$2"
      shift 2
      ;;
    --threshold-mutation)
      THRESHOLD_MUTATION="$2"
      shift 2
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

cmd=(python3 -m benchmark.reporter --results-root "$RESULTS_ROOT" --formats "$FORMATS" --chart-style "$CHART_STYLE" --top-n-errors "$TOP_N_ERRORS" --threshold-compile "$THRESHOLD_COMPILE" --threshold-test "$THRESHOLD_TEST" --threshold-line "$THRESHOLD_LINE" --threshold-branch "$THRESHOLD_BRANCH" --threshold-mutation "$THRESHOLD_MUTATION")

if [[ -n "$EVALUATOR_SUMMARY" ]]; then
  cmd+=(--evaluator-summary "$EVALUATOR_SUMMARY")
fi

if [[ -n "$OUTPUT_DIR" ]]; then
  cmd+=(--output-dir "$OUTPUT_DIR")
fi

echo "[gen_report] running: ${cmd[*]}"
"${cmd[@]}"
