from __future__ import annotations

import argparse
import json
import logging

from .runner import Runner


def _split_csv(raw: str | None) -> list[str] | None:
    if not raw:
        return None
    values = [x.strip() for x in raw.split(",")]
    return [x for x in values if x] or None


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="ut-bench runner: send dataset samples to AI models"
    )
    parser.add_argument(
        "--config",
        default="benchmark/config/models.yaml",
        help="Path to models.yaml",
    )
    parser.add_argument(
        "--dataset-root",
        default="dataset",
        help="Dataset directory",
    )
    parser.add_argument(
        "--results-root",
        default="results",
        help="Output root directory",
    )
    parser.add_argument(
        "--model",
        default=None,
        help="Comma-separated model names, e.g. deepseek,glm",
    )
    parser.add_argument(
        "--lang",
        default=None,
        help="Comma-separated languages, e.g. python,java,go",
    )
    parser.add_argument(
        "--max-samples",
        type=int,
        default=None,
        help="Limit samples per language directory",
    )
    parser.add_argument(
        "--sample-glob",
        default=None,
        help="Glob pattern under each language dir, e.g. boundary/*.py",
    )
    parser.add_argument(
        "--retries",
        type=int,
        default=3,
        help="Max retries for API calls",
    )
    parser.add_argument(
        "--backoff-seconds",
        type=float,
        default=2.0,
        help="Base exponential backoff seconds",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Do not call API; generate placeholder tests",
    )
    parser.add_argument(
        "--max-workers",
        type=int,
        default=None,
        help="Override total worker threads for concurrent execution",
    )
    parser.add_argument(
        "--no-resume",
        action="store_true",
        help="Disable checkpoint resume; run all tasks in current scope",
    )
    parser.add_argument(
        "--reset-checkpoint",
        action="store_true",
        help="Delete scope checkpoint before running",
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Enable debug logging",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    logging.basicConfig(
        level=logging.DEBUG if args.verbose else logging.INFO,
        format="[%(asctime)s] [%(levelname)s] %(message)s",
    )

    runner = Runner(
        config_path=args.config,
        dataset_root=args.dataset_root,
        results_root=args.results_root,
        retries=args.retries,
        backoff_seconds=args.backoff_seconds,
    )

    summary = runner.run(
        model_names=_split_csv(args.model),
        languages=_split_csv(args.lang),
        max_samples=args.max_samples,
        sample_glob=args.sample_glob,
        dry_run=args.dry_run,
        resume=not args.no_resume,
        reset_checkpoint=args.reset_checkpoint,
        max_workers=args.max_workers,
    )

    print(json.dumps(summary, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
