import json
from collections import Counter, defaultdict
from pathlib import Path


CATEGORIES = ["simple_function", "boundary", "complex_dependency", "interface_mock"]
DEFAULT_SELECTION_ORDER = ["interface_mock", "complex_dependency", "boundary", "simple_function"]


def ensure_dir(path: Path) -> None:
    path.mkdir(parents=True, exist_ok=True)


def write_json(path: Path, payload: dict) -> None:
    ensure_dir(path.parent)
    path.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")


def effective_line_count(code: str) -> int:
    return sum(1 for line in code.splitlines() if line.strip())


def select_balanced_records(
    candidates: list[dict],
    *,
    categories: list[str] | None = None,
    target_per_category: int = 50,
    selection_order: list[str] | None = None,
    unique_key: str = "id",
) -> tuple[list[dict], dict[str, list[dict]]]:
    categories = categories or CATEGORIES
    selection_order = selection_order or DEFAULT_SELECTION_ORDER

    per_category: dict[str, list[dict]] = defaultdict(list)
    for candidate in candidates:
        for category in categories:
            if candidate["category_eligibility"].get(category):
                per_category[category].append(candidate)

    for category in categories:
        per_category[category].sort(
            key=lambda item: (
                -item["category_scores"][category],
                -item.get("complexity_score", 0),
                -item.get("effective_lines", 0),
                str(item.get(unique_key, "")),
            )
        )

    selected = []
    used_keys = set()
    for category in selection_order:
        picked = 0
        for candidate in per_category.get(category, []):
            if picked >= target_per_category:
                break
            key = candidate[unique_key]
            if key in used_keys:
                continue
            selected.append({**candidate, "category": category})
            used_keys.add(key)
            picked += 1

    return selected, per_category


def category_counts(records: list[dict]) -> dict[str, int]:
    return dict(Counter(item["category"] for item in records))


def build_category_stats(records: list[dict], categories: list[str] | None = None) -> dict[str, dict]:
    categories = categories or CATEGORIES
    stats = {}
    for category in categories:
        bucket = [item for item in records if item["category"] == category]
        lengths = [len(item["code"]) for item in bucket]
        stats[category] = {
            "count": len(bucket),
            "avg_effective_lines": round(
                sum(item.get("effective_lines", effective_line_count(item["code"])) for item in bucket) / len(bucket), 2
            )
            if bucket
            else 0,
            "min_length": min(lengths) if lengths else 0,
            "max_length": max(lengths) if lengths else 0,
        }
    return stats
