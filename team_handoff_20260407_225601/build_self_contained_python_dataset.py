import ast
import hashlib
import json
import warnings
from collections import Counter, defaultdict

from build_balanced_python_dataset import CATEGORIES, TARGET_PER_CATEGORY, category_eligibility, category_scores
from python_dataset_precheck import analyze_self_contained_python
from select_complex_python_dataset import BASE_DIR, ComplexityVisitor, SOURCE_FILE, count_effective_lines, strip_comments_and_docstrings


OUTPUT_DIR = BASE_DIR / "data" / "python_code_files_self_contained"
OUTPUT_JSON = BASE_DIR / "data" / "python_ut_dataset" / "python_dataset_self_contained.json"


def sha1_text(text: str) -> str:
    return hashlib.sha1(text.encode("utf-8", errors="ignore")).hexdigest()


def build_relaxed_record(task_id: str, original_code: str):
    cleaned_code = strip_comments_and_docstrings(original_code).strip() + "\n"
    try:
        with warnings.catch_warnings():
            warnings.simplefilter("ignore", SyntaxWarning)
            tree = ast.parse(cleaned_code)
    except SyntaxError:
        return None, None

    visitor = ComplexityVisitor()
    visitor.visit(tree)
    effective_lines = count_effective_lines(cleaned_code)
    has_nested_loop = visitor.nested_loop_count > 0
    has_recursion = bool(visitor.recursive_functions)
    has_try_except = visitor.try_count > 0
    has_class = visitor.class_count > 0
    has_stateful_oop = has_class and visitor.self_mutations > 0
    has_multiple_branches = visitor.branch_count >= 3

    if effective_lines < 8:
        return None, None
    if visitor.cyclomatic < 2:
        return None, None

    score = (
        visitor.cyclomatic * 3
        + effective_lines * 0.35
        + visitor.max_depth * 4
        + visitor.branch_count * 2
        + visitor.try_count * 4
        + visitor.nested_loop_count * 5
        + len(visitor.recursive_functions) * 6
        + visitor.self_mutations * 3
        + visitor.collection_mutations
    )

    return (
        {
            "task_id": task_id,
            "code": cleaned_code,
            "effective_lines": effective_lines,
            "cyclomatic_complexity": visitor.cyclomatic,
            "complexity_score": round(score, 2),
            "max_nesting_depth": visitor.max_depth,
            "branch_count": visitor.branch_count,
            "try_count": visitor.try_count,
            "class_count": visitor.class_count,
            "self_mutations": visitor.self_mutations,
            "collection_mutations": visitor.collection_mutations,
            "has_nested_loop": has_nested_loop,
            "has_recursion": has_recursion,
            "has_try_except": has_try_except,
            "has_class": has_class,
            "has_stateful_oop": has_stateful_oop,
            "has_multiple_branches": has_multiple_branches,
        },
        tree,
    )


def iter_candidates():
    seen_hashes = set()
    with SOURCE_FILE.open("r", encoding="utf-8") as handle:
        for line in handle:
            if not line.strip():
                continue
            item = json.loads(line)
            task_id = item.get("task_id", "unknown")

            prompt = item.get("complete_prompt", item.get("code_prompt", ""))
            solution = item.get("canonical_solution", "")
            code = f"{prompt}{solution}"
            record, tree = build_relaxed_record(task_id, code)
            if record is None:
                continue

            code_hash = sha1_text(record["code"])
            if code_hash in seen_hashes:
                continue
            seen_hashes.add(code_hash)

            self_contained = analyze_self_contained_python(tree)
            if not self_contained.is_self_contained:
                continue

            test_ref = item.get("test", "")
            text = f"{code}\n{test_ref}\n{item.get('libs', '')}".lower()
            scores = category_scores(record, text)
            eligibility = category_eligibility(record, text)
            valid_scores = {category: value for category, value in scores.items() if eligibility[category]}
            if not valid_scores:
                continue

            yield {
                **record,
                "source": "BigCodeBench",
                "test_ref": test_ref,
                "libs": item.get("libs", ""),
                "category_scores": scores,
                "category_eligibility": eligibility,
                "suggested_category": max(valid_scores, key=valid_scores.get),
                "self_contained_check": self_contained.to_dict(),
            }


def select_balanced_records(candidates):
    per_category = defaultdict(list)
    for candidate in candidates:
        for category in CATEGORIES:
            if candidate["category_eligibility"][category]:
                per_category[category].append(candidate)

    for category in CATEGORIES:
        per_category[category].sort(
            key=lambda item: (
                -item["category_scores"][category],
                -item["complexity_score"],
                -item["cyclomatic_complexity"],
                -item["effective_lines"],
                item["task_id"],
            )
        )

    selection_order = ["interface_mock", "boundary", "complex_dependency", "simple_function"]
    candidates_sorted = sorted(
        candidates,
        key=lambda item: (
            -max(item["category_scores"][category] for category in CATEGORIES if item["category_eligibility"][category]),
            sum(item["category_eligibility"].values()),
            -item["complexity_score"],
            -item["cyclomatic_complexity"],
            -item["effective_lines"],
            item["task_id"],
        ),
    )

    category_slots = {
        category: [(category, slot_index) for slot_index in range(TARGET_PER_CATEGORY)]
        for category in selection_order
    }
    slot_to_candidate: dict[tuple[str, int], int] = {}

    def iter_candidate_slots(candidate: dict):
        eligible_categories = [
            category for category in selection_order if candidate["category_eligibility"][category]
        ]
        eligible_categories.sort(
            key=lambda category: (
                -candidate["category_scores"][category],
                selection_order.index(category),
            )
        )
        for category in eligible_categories:
            yield from category_slots[category]

    def try_assign(candidate_index: int, seen_slots: set[tuple[str, int]]) -> bool:
        candidate = candidates_sorted[candidate_index]
        for slot in iter_candidate_slots(candidate):
            if slot in seen_slots:
                continue
            seen_slots.add(slot)
            current_owner = slot_to_candidate.get(slot)
            if current_owner is None or try_assign(current_owner, seen_slots):
                slot_to_candidate[slot] = candidate_index
                return True
        return False

    for candidate_index in range(len(candidates_sorted)):
        try_assign(candidate_index, set())

    selected = []
    for (category, _), candidate_index in slot_to_candidate.items():
        selected.append({**candidates_sorted[candidate_index], "category": category})

    counts = Counter(item["category"] for item in selected)
    missing = [category for category in CATEGORIES if counts[category] < TARGET_PER_CATEGORY]
    if missing:
        raise RuntimeError(
            "Unable to satisfy self-contained dataset targets: "
            + ", ".join(f"{category}={counts[category]}/{TARGET_PER_CATEGORY}" for category in missing)
        )

    return selected, per_category


def write_outputs(selected_records):
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    category_stats = {}
    samples = []

    for category in CATEGORIES:
        category_dir = OUTPUT_DIR / category
        category_dir.mkdir(parents=True, exist_ok=True)
        category_records = [item for item in selected_records if item["category"] == category]

        lengths = [len(item["code"]) for item in category_records]
        ccs = [item["cyclomatic_complexity"] for item in category_records]
        category_stats[category] = {
            "count": len(category_records),
            "avg_effective_lines": round(
                sum(item["effective_lines"] for item in category_records) / len(category_records), 2
            ),
            "avg_cyclomatic_complexity": round(sum(ccs) / len(ccs), 2),
            "min_length": min(lengths),
            "max_length": max(lengths),
        }

        for index, record in enumerate(category_records):
            file_name = f"{category}_{index:03d}.py"
            file_path = category_dir / file_name
            file_path.write_text(record["code"], encoding="utf-8")
            samples.append(
                {
                    "id": f"python_{category}_{index}",
                    "language": "python",
                    "category": category,
                    "source": record["source"],
                    "task_id": record["task_id"],
                    "code_file": str(file_path.relative_to(BASE_DIR)).replace("\\", "/"),
                    "code": record["code"],
                    "test_ref": record["test_ref"],
                    "metrics": {
                        "effective_lines": record["effective_lines"],
                        "cyclomatic_complexity": record["cyclomatic_complexity"],
                        "complexity_score": record["complexity_score"],
                        "max_nesting_depth": record["max_nesting_depth"],
                        "branch_count": record["branch_count"],
                        "try_count": record["try_count"],
                    },
                    "features": {
                        "has_nested_loop": record["has_nested_loop"],
                        "has_recursion": record["has_recursion"],
                        "has_try_except": record["has_try_except"],
                        "has_class": record["has_class"],
                        "has_stateful_oop": record["has_stateful_oop"],
                        "has_multiple_branches": record["has_multiple_branches"],
                    },
                    "category_scores": record["category_scores"],
                    "self_contained_check": record["self_contained_check"],
                }
            )

    payload = {
        "metadata": {
            "version": "Python_UT_Self_Contained_v1",
            "description": "Balanced Python dataset for unit-test evaluation with single-file, stdlib-only imports",
            "source_file": str(SOURCE_FILE.relative_to(BASE_DIR)).replace("\\", "/"),
            "selection_strategy": {
                "per_category": TARGET_PER_CATEGORY,
                "min_effective_lines": 8,
                "min_cyclomatic_complexity": 2,
                "source": "BigCodeBench",
                "deduplicate_by": "cleaned_source_sha1",
                "require_self_contained_single_file": True,
                "allow_stdlib_imports_only": True,
            },
            "categories": {
                "simple_function": "Algorithmic or transformation-heavy functions with limited external dependency",
                "boundary": "Validation, error handling, and edge-case oriented logic",
                "complex_dependency": "Code involving file, network, data, async, process, or multi-step dependency logic",
                "interface_mock": "Code whose behavior is naturally evaluated with mocks, stubs, or patched collaborators",
            },
        },
        "statistics": {
            "total": len(samples),
            "by_category": category_stats,
            "overall": {
                "avg_effective_lines": round(
                    sum(item["metrics"]["effective_lines"] for item in samples) / len(samples), 2
                ),
                "avg_cyclomatic_complexity": round(
                    sum(item["metrics"]["cyclomatic_complexity"] for item in samples) / len(samples), 2
                ),
            },
        },
        "samples": samples,
    }
    OUTPUT_JSON.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")


def main():
    candidates = list(iter_candidates())
    selected, per_category = select_balanced_records(candidates)
    counts = Counter(item["category"] for item in selected)

    write_outputs(selected)
    print("=== Self-contained Python dataset ===")
    print(f"Candidates: {len(candidates)}")
    for category in CATEGORIES:
        print(f"{category}: selected {counts[category]}, candidate pool {len(per_category[category])}")
    print(f"Output directory: {OUTPUT_DIR}")
    print(f"Output json: {OUTPUT_JSON}")


if __name__ == "__main__":
    main()
