import ast
import hashlib
import json
import warnings
from collections import Counter, defaultdict
from pathlib import Path

from build_realworld_module_dataset import (
    BASE_DIR,
    BOUNDARY_KEYWORDS,
    CATEGORIES,
    DEPENDENCY_IMPORT_HINTS,
    DEPENDENCY_KEYWORDS,
    EXCLUDE_NAMES,
    EXCLUDE_PARTS,
    INTERFACE_KEYWORDS,
    MIN_CYCLOMATIC,
    MIN_EFFECTIVE_LINES,
    MAX_EFFECTIVE_LINES,
    SOURCE_ROOTS,
    classify_scores,
    infer_package,
)
from dataset_builder_utils import ensure_dir, write_json
from python_dataset_precheck import analyze_self_contained_python
from select_complex_python_dataset import ComplexityVisitor, count_effective_lines, strip_comments_and_docstrings


OUTPUT_DIR = BASE_DIR / "data" / "python_code_files_module_level"
OUTPUT_JSON = BASE_DIR / "data" / "python_ut_dataset" / "python_dataset_module_level.json"
TARGET_PER_CATEGORY = 50


def sha1_text(text: str) -> str:
    return hashlib.sha1(text.encode("utf-8", errors="ignore")).hexdigest()


def iter_candidates():
    seen_hashes = set()
    for root in SOURCE_ROOTS:
        if not root.exists():
            continue
        for path in root.rglob("*.py"):
            if path.name in EXCLUDE_NAMES:
                continue
            if any(part in EXCLUDE_PARTS for part in path.parts):
                continue
            try:
                original = path.read_text(encoding="utf-8")
                cleaned = strip_comments_and_docstrings(original).strip() + "\n"
                with warnings.catch_warnings():
                    warnings.simplefilter("ignore", SyntaxWarning)
                    tree = ast.parse(cleaned)
            except Exception:
                continue

            visitor = ComplexityVisitor()
            visitor.visit(tree)
            effective_lines = count_effective_lines(cleaned)
            funcs = sum(isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)) for node in ast.walk(tree))
            classes = sum(isinstance(node, ast.ClassDef) for node in ast.walk(tree))
            if effective_lines < MIN_EFFECTIVE_LINES or effective_lines > MAX_EFFECTIVE_LINES:
                continue
            if visitor.cyclomatic < MIN_CYCLOMATIC:
                continue
            if funcs + classes < 3:
                continue

            text_hash = sha1_text(cleaned)
            if text_hash in seen_hashes:
                continue
            seen_hashes.add(text_hash)

            package = infer_package(path)
            self_contained = analyze_self_contained_python(tree, package=package)
            record = {
                "source_path": str(path.relative_to(BASE_DIR)).replace("\\", "/"),
                "package": package,
                "original_code": original,
                "code": cleaned,
                "effective_lines": effective_lines,
                "cyclomatic_complexity": visitor.cyclomatic,
                "max_nesting_depth": visitor.max_depth,
                "branch_count": visitor.branch_count,
                "try_count": visitor.try_count,
                "class_count": classes,
                "function_count": funcs,
                "complexity_score": (
                    visitor.cyclomatic * 3
                    + effective_lines * 0.3
                    + visitor.max_depth * 4
                    + visitor.branch_count
                    + visitor.try_count * 3
                ),
                "self_contained_check": self_contained.to_dict(),
            }
            eligibility, scores = classify_scores(
                record=record,
                path_text=record["source_path"],
                code_text=cleaned,
                funcs=funcs,
                classes=classes,
            )
            if not any(eligibility.values()):
                continue
            record["category_eligibility"] = eligibility
            record["category_scores"] = scores
            yield record


def select_records(candidates):
    by_category = defaultdict(list)
    for candidate in candidates:
        for category in CATEGORIES:
            if candidate["category_eligibility"][category]:
                by_category[category].append(candidate)

    for category in CATEGORIES:
        by_category[category].sort(
            key=lambda item: (
                -item["category_scores"][category],
                -item["complexity_score"],
                -item["cyclomatic_complexity"],
                -item["effective_lines"],
                item["source_path"],
            )
        )

    selected = []
    used_paths = set()
    for category in ["interface_mock", "complex_dependency", "boundary", "simple_function"]:
        picked = 0
        for candidate in by_category[category]:
            if picked >= TARGET_PER_CATEGORY:
                break
            if candidate["source_path"] in used_paths:
                continue
            selected.append({**candidate, "category": category})
            used_paths.add(candidate["source_path"])
            picked += 1

    counts = Counter(item["category"] for item in selected)
    missing = [category for category in CATEGORIES if counts[category] < TARGET_PER_CATEGORY]
    if missing:
        raise RuntimeError(
            "Not enough Python module-level candidates: "
            + ", ".join(f"{category}={counts[category]}/{TARGET_PER_CATEGORY}" for category in missing)
        )
    return selected, by_category


def write_outputs(selected):
    ensure_dir(OUTPUT_DIR)
    samples = []
    category_stats = {}

    for category in CATEGORIES:
        category_dir = OUTPUT_DIR / category
        ensure_dir(category_dir)
        records = [item for item in selected if item["category"] == category]
        lengths = [len(item["code"]) for item in records]
        category_stats[category] = {
            "count": len(records),
            "avg_effective_lines": round(sum(item["effective_lines"] for item in records) / len(records), 2),
            "avg_cyclomatic_complexity": round(
                sum(item["cyclomatic_complexity"] for item in records) / len(records), 2
            ),
            "min_length": min(lengths),
            "max_length": max(lengths),
        }

        for index, record in enumerate(records):
            file_name = f"{category}_{index:03d}.py"
            out_path = category_dir / file_name
            out_path.write_text(record["code"], encoding="utf-8")
            samples.append(
                {
                    "id": f"python_{category}_{index}",
                    "language": "python",
                    "category": category,
                    "source": "project_module",
                    "package": record["package"],
                    "source_path": record["source_path"],
                    "code_file": str(out_path.relative_to(BASE_DIR)).replace("\\", "/"),
                    "code": record["code"],
                    "original_length": len(record["original_code"]),
                    "metrics": {
                        "effective_lines": record["effective_lines"],
                        "cyclomatic_complexity": record["cyclomatic_complexity"],
                        "complexity_score": round(record["complexity_score"], 2),
                        "max_nesting_depth": record["max_nesting_depth"],
                        "branch_count": record["branch_count"],
                        "try_count": record["try_count"],
                        "function_count": record["function_count"],
                        "class_count": record["class_count"],
                    },
                    "category_scores": record["category_scores"],
                    "self_contained_check": record["self_contained_check"],
                }
            )

    payload = {
        "metadata": {
            "version": "Python_Module_Level_v1",
            "description": "Balanced project/module-level Python dataset using real modules without self-contained filtering",
            "source_roots": [str(path.relative_to(BASE_DIR)).replace("\\", "/") for path in SOURCE_ROOTS if path.exists()],
            "selection_strategy": {
                "per_category": TARGET_PER_CATEGORY,
                "min_effective_lines": MIN_EFFECTIVE_LINES,
                "max_effective_lines": MAX_EFFECTIVE_LINES,
                "min_cyclomatic_complexity": MIN_CYCLOMATIC,
                "exclude_tests": True,
                "deduplicate_by": "cleaned_source_sha1",
                "require_self_contained_single_file": False,
                "requires_project_context": True,
            },
            "categories": {
                "simple_function": "Function-heavy modules with mostly algorithmic or transformation logic",
                "boundary": "Modules rich in validation, parsing, normalization, and error handling",
                "complex_dependency": "Modules with file, network, data, async, process, or multi-library dependencies",
                "interface_mock": "Class-oriented integration surfaces that are natural targets for mocks, stubs, or patched collaborators",
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
                "packages": dict(Counter(item["package"] for item in samples)),
                "self_contained_failures": dict(
                    Counter(
                        issue["issue_type"]
                        for item in samples
                        for issue in item["self_contained_check"]["issues"]
                    )
                ),
            },
        },
        "samples": samples,
    }
    write_json(OUTPUT_JSON, payload)


def main():
    candidates = list(iter_candidates())
    selected, by_category = select_records(candidates)
    write_outputs(selected)
    counts = Counter(item["category"] for item in selected)
    print("=== Python module-level dataset ===")
    print(f"Candidates: {len(candidates)}")
    for category in CATEGORIES:
        print(f"{category}: selected {counts[category]}, candidate pool {len(by_category[category])}")
    print(f"Output directory: {OUTPUT_DIR}")
    print(f"Output JSON: {OUTPUT_JSON}")


if __name__ == "__main__":
    main()
