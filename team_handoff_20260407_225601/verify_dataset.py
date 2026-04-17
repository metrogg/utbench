import argparse
import ast
import json
from collections import Counter
from pathlib import Path

from python_dataset_precheck import analyze_self_contained_python


DEFAULT_DATASET = Path("data/python_ut_dataset/python_dataset_self_contained.json")


def load_dataset(dataset_path: Path) -> dict:
    return json.loads(dataset_path.read_text(encoding="utf-8"))


def summarize_dataset(data: dict) -> None:
    print("=== Python Dataset Statistics ===")
    print(f"Total samples: {data['statistics']['total']}")
    print()

    for category, stats in data["statistics"]["by_category"].items():
        print(f"{category}:")
        print(f"  Count: {stats['count']}")
        print(f"  Min length: {stats['min_length']} chars")
        print(f"  Max length: {stats['max_length']} chars")
        if "avg_effective_lines" in stats:
            print(f"  Avg effective lines: {stats['avg_effective_lines']}")
        if "avg_cyclomatic_complexity" in stats:
            print(f"  Avg cyclomatic complexity: {stats['avg_cyclomatic_complexity']}")
        print()


def validate_self_contained_samples(data: dict) -> tuple[Counter, list[tuple[str, list[str]]]]:
    issue_counts: Counter = Counter()
    failures: list[tuple[str, list[str]]] = []

    for sample in data["samples"]:
        tree = ast.parse(sample["code"])
        result = analyze_self_contained_python(tree)
        if result.is_self_contained:
            continue
        issues = sorted({issue.issue_type for issue in result.issues})
        for issue in issues:
            issue_counts[issue] += 1
        failures.append((sample["id"], issues))

    return issue_counts, failures


def main() -> int:
    parser = argparse.ArgumentParser(description="Validate Python dataset executability.")
    parser.add_argument(
        "dataset",
        nargs="?",
        default=str(DEFAULT_DATASET),
        help="Path to a dataset json file.",
    )
    args = parser.parse_args()

    dataset_path = Path(args.dataset)
    data = load_dataset(dataset_path)
    summarize_dataset(data)

    issue_counts, failures = validate_self_contained_samples(data)
    print("=== Self-contained Check ===")
    print(f"Dataset file: {dataset_path}")
    print(f"Samples with import issues: {len(failures)}")
    print(f"Issue counts: {dict(issue_counts)}")

    if failures:
        print("Examples:")
        for sample_id, issues in failures[:10]:
            print(f"  {sample_id}: {', '.join(issues)}")
        return 1

    print("All samples passed the single-file stdlib-only import check.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
