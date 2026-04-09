from __future__ import annotations

import argparse
import json

from .pipeline import Evaluator


def _split_csv(raw: str | None) -> list[str] | None:
    if not raw:
        return None
    values = [item.strip() for item in raw.split(",")]
    return [item for item in values if item] or None


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="ut-bench evaluator")
    parser.add_argument("--results-root", default="results")
    parser.add_argument("--model", default=None, help="Comma-separated model names")
    parser.add_argument("--lang", default=None, help="Comma-separated languages")
    parser.add_argument("--output-dir", default=None)
    parser.add_argument(
        "--all-history",
        action="store_true",
        help="Evaluate all historical test files per sample instead of latest only",
    )
    parser.add_argument(
        "--only-self-contained",
        action="store_true",
        help="Evaluate only Python samples in self-contained allowlist",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    evaluator = Evaluator(results_root=args.results_root)
    payload = evaluator.run(
        models=_split_csv(args.model),
        languages=_split_csv(args.lang),
        output_dir=args.output_dir,
        latest_only=not args.all_history,
        only_self_contained=args.only_self_contained,
    )
    print(json.dumps(payload["summary"], ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
