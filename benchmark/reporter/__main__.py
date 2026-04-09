from __future__ import annotations

import argparse
import json

from .pipeline import Reporter


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="ut-bench reporter")
    parser.add_argument("--results-root", default="results")
    parser.add_argument(
        "--evaluator-summary",
        default=None,
        help="Path to evaluator_summary_*.json; default uses latest under results-root",
    )
    parser.add_argument(
        "--output-dir",
        default=None,
        help="Output directory for reporter artifacts; default is <results-root>/reports",
    )
    parser.add_argument("--top-n-errors", type=int, default=20)
    parser.add_argument(
        "--formats",
        default="json,csv,html",
        help="Comma-separated output formats: json,csv,html",
    )
    parser.add_argument(
        "--chart-style",
        default="teal",
        help="Chart style palette: teal, orange, blue",
    )
    parser.add_argument("--threshold-compile", type=float, default=1.0)
    parser.add_argument("--threshold-test", type=float, default=0.7)
    parser.add_argument("--threshold-line", type=float, default=0.7)
    parser.add_argument("--threshold-branch", type=float, default=0.6)
    parser.add_argument("--threshold-mutation", type=float, default=0.85)
    return parser.parse_args()


def _split_csv(raw: str | None) -> list[str] | None:
    if not raw:
        return None
    values = [item.strip() for item in raw.split(",")]
    return [item for item in values if item] or None


def main() -> int:
    args = parse_args()
    reporter = Reporter(results_root=args.results_root)
    payload = reporter.run(
        evaluator_summary=args.evaluator_summary,
        output_dir=args.output_dir,
        top_n_errors=args.top_n_errors,
        formats=_split_csv(args.formats),
        chart_style=args.chart_style,
        thresholds={
            "compile_pass_rate": args.threshold_compile,
            "test_pass_rate": args.threshold_test,
            "line_coverage": args.threshold_line,
            "branch_coverage": args.threshold_branch,
            "mutation_score": args.threshold_mutation,
        },
    )
    print(json.dumps(payload["summary"], ensure_ascii=False, indent=2))
    print(json.dumps(payload["artifacts"], ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
