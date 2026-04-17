import json
import re
from collections import Counter
from pathlib import Path

from dataset_builder_utils import CATEGORIES, build_category_stats, ensure_dir, select_balanced_records, write_json
from sync_external_self_contained_sources import OUTPUT_DIR as EXTERNAL_SOURCE_DIR
from sync_external_self_contained_sources import main as sync_external_sources


BASE_DIR = Path(__file__).resolve().parent
MODULE_SOURCE_ROOT = BASE_DIR / "data" / "bcb" / "code" / "hm-dianping" / "src" / "main" / "java"
MODULE_OUTPUT_DIR = BASE_DIR / "data" / "java_code_files_module_level"
MODULE_OUTPUT_JSON = BASE_DIR / "data" / "java_ut_dataset" / "java_dataset_module_level.json"
SELF_CONTAINED_OUTPUT_DIR = BASE_DIR / "data" / "java_code_files_self_contained"
SELF_CONTAINED_OUTPUT_JSON = BASE_DIR / "data" / "java_ut_dataset" / "java_dataset_self_contained.json"
TARGET_PER_CATEGORY = 50


def effective_line_count(code: str) -> int:
    return sum(1 for line in code.splitlines() if line.strip())


def java_module_category_scores(path_text: str, code: str) -> dict[str, float]:
    low = f"{path_text}\n{code}".lower()
    simple = (
        2 * sum(token in low for token in ["return ", "if (", "for ("])
        + 4 * sum(token in low for token in ["/dto/", "/entity/", "constants", "pattern", "regex"])
        - 3 * sum(token in low for token in ["@service", "@restcontroller", "redis", "session", "mapper"])
    )
    boundary = (
        4 * sum(token in low for token in ["regex", "pattern", "invalid", "null", "exception", "advice", "error"])
        + low.count("throw ")
        + low.count("if (")
    )
    complex_dependency = (
        5 * sum(token in low for token in ["redis", "session", "mapper", "serviceimpl", "mybatis", "stream", "cache"])
        + 4 * sum(token in low for token in ["/service/impl/", "/config/", "/controller/"])
        + low.count("@autowired")
    )
    interface_mock = (
        5 * sum(token in low for token in ["interface ", "@service", "@restcontroller", "mapper", "serviceimpl"])
        + 3 * sum(token in low for token in ["/service/", "/mapper/", "/controller/"])
        + low.count("@override")
    )
    return {
        "simple_function": simple,
        "boundary": boundary,
        "complex_dependency": complex_dependency,
        "interface_mock": interface_mock,
    }


def infer_package_name(code: str) -> str:
    match = re.search(r"^package\s+([^;]+);", code, re.M)
    return match.group(1) if match else ""


def build_java_module_dataset() -> None:
    ensure_dir(MODULE_OUTPUT_DIR)
    samples = []
    buckets = {category: [] for category in CATEGORIES}

    for path in MODULE_SOURCE_ROOT.rglob("*.java"):
        code = path.read_text(encoding="utf-8", errors="ignore")
        effective_lines = effective_line_count(code)
        if effective_lines < 8:
            continue
        rel_path = str(path.relative_to(BASE_DIR)).replace("\\", "/")
        scores = java_module_category_scores(rel_path, code)
        category = max(scores, key=scores.get)
        buckets[category].append(
            {
                "language": "java",
                "source": "hm-dianping",
                "project": "hm-dianping",
                "source_path": rel_path,
                "package_name": infer_package_name(code),
                "code": code,
                "effective_lines": effective_lines,
                "category_scores": scores,
            }
        )

    for category in CATEGORIES:
        bucket_dir = MODULE_OUTPUT_DIR / category
        ensure_dir(bucket_dir)
        records = sorted(
            buckets[category],
            key=lambda item: (-item["category_scores"][category], -item["effective_lines"], item["source_path"]),
        )
        for index, record in enumerate(records):
            file_path = bucket_dir / f"{category}_{index:03d}.java"
            file_path.write_text(record["code"], encoding="utf-8")
            samples.append(
                {
                    "id": f"java_{category}_{index}",
                    "language": "java",
                    "category": category,
                    "source": record["source"],
                    "project": record["project"],
                    "package_name": record["package_name"],
                    "source_path": record["source_path"],
                    "code_file": str(file_path.relative_to(BASE_DIR)).replace("\\", "/"),
                    "code": record["code"],
                    "effective_lines": record["effective_lines"],
                    "category_scores": record["category_scores"],
                }
            )

    write_json(
        MODULE_OUTPUT_JSON,
        {
            "metadata": {
                "version": "Java_Module_Level_v1",
                "language": "java",
                "source_root": str(MODULE_SOURCE_ROOT.relative_to(BASE_DIR)).replace("\\", "/"),
                "description": "Java project/module-level dataset built from the hm-dianping project source tree",
                "selection_strategy": {
                    "balanced_categories": False,
                    "requires_project_context": True,
                    "granularity": "project_module_file",
                },
            },
            "statistics": {
                "total": len(samples),
                "by_category": build_category_stats(samples),
                "by_source": {"hm-dianping": len(samples)},
            },
            "samples": samples,
        },
    )


def load_jsonl(path: Path) -> list[dict]:
    with path.open("r", encoding="utf-8") as handle:
        return [json.loads(line) for line in handle if line.strip()]


def java_self_contained_scores(record: dict) -> dict[str, float]:
    code = record["code"]
    question = record.get("question", "")
    difficulty = record.get("difficulty", "")
    low = f"{question}\n{code}\n{record.get('test_ref', '')}".lower()
    class_count = len(re.findall(r"\bclass\s+\w+", code))
    method_count = len(re.findall(r"\b(public|private|protected)\b", code))
    source_bias = 4 if record["source"] == "HumanEval-X" else 0
    oop_bias = 4 if record["source"] == "AutoCodeBenchmark" else 0

    return {
        "simple_function": (
            3 * sum(token in low for token in ["return ", "for (", "while (", "math.", "arraylist", "list<", "sort"])
            + source_bias
            + (2 if difficulty == "easy" else 0)
            - class_count
        ),
        "boundary": (
            4 * sum(token in low for token in ["throw", "exception", "invalid", "null", "empty", "range", "regex"])
            + low.count("if (")
            + source_bias
            + (2 if difficulty == "medium" else 0)
        ),
        "complex_dependency": (
            4
            * sum(
                token in low
                for token in ["hashmap", "treemap", "tree", "graph", "matrix", "queue", "stack", "cache", "date", "time"]
            )
            + method_count
            + oop_bias
            + (3 if difficulty == "hard" else 0)
        ),
        "interface_mock": (
            5 * sum(token in low for token in ["interface", "manager", "system", "service", "analyzer", "account", "user"])
            + 3 * class_count
            + oop_bias
            + (3 if difficulty == "hard" else 0)
        ),
    }


def is_self_contained_java(code: str) -> tuple[bool, list[str]]:
    issues = []
    imports = re.findall(r"^import\s+([^;]+);", code, re.M)
    for imp in imports:
        if not (imp.startswith("java.") or imp.startswith("javax.")):
            issues.append(f"non_std_import:{imp}")
    if re.search(r"^package\s+", code, re.M):
        issues.append("package_declaration")
    if re.search(r"^\s*import\s+static\s+", code, re.M):
        issues.append("static_import")
    if re.search(r"\b(org|com|net)\.[A-Za-z0-9_.]+", code):
        issues.append("third_party_reference")
    return not issues, issues


def load_java_self_contained_candidates() -> list[dict]:
    sync_external_sources()
    candidates = []

    for item in load_jsonl(EXTERNAL_SOURCE_DIR / "humanevalx_java.jsonl"):
        code = (item["prompt"] + item["canonical_solution"]).strip() + "\n"
        ok, issues = is_self_contained_java(code)
        if not ok:
            continue
        scores = java_self_contained_scores(
            {
                "source": "HumanEval-X",
                "code": code,
                "test_ref": item.get("test", ""),
                "question": item.get("text", "") or item.get("prompt", ""),
                "difficulty": "benchmark",
            }
        )
        candidates.append(
            {
                "id": f"humanevalx_{item['task_id'].replace('/', '_').lower()}",
                "source": "HumanEval-X",
                "source_record_id": item["task_id"],
                "source_path": "data/external_sources/java_cpp_self_contained/humanevalx_java.jsonl",
                "code": code,
                "test_ref": item.get("test", ""),
                "question": item.get("text", "") or item.get("prompt", ""),
                "difficulty": "benchmark",
                "effective_lines": effective_line_count(code),
                "complexity_score": len(code) / 80 + code.count("if (") + code.count("for (") + 2 * code.count("class "),
                "category_scores": scores,
                "category_eligibility": {
                    "simple_function": scores["simple_function"] >= 6,
                    "boundary": scores["boundary"] >= 5,
                    "complex_dependency": scores["complex_dependency"] >= 5,
                    "interface_mock": scores["interface_mock"] >= 3,
                },
            }
        )

    acb_index = 0
    for item in load_jsonl(EXTERNAL_SOURCE_DIR / "autocodebenchmark_java.jsonl"):
        code = item["canonical_solution"].strip()
        if not code:
            continue
        code = code + "\n"
        ok, issues = is_self_contained_java(code)
        if not ok:
            continue
        test_ref = item.get("full_test_func") or item.get("demo_test_func") or ""
        scores = java_self_contained_scores(
            {
                "source": "AutoCodeBenchmark",
                "code": code,
                "test_ref": test_ref,
                "question": item.get("question", ""),
                "difficulty": item.get("difficulty", ""),
            }
        )
        acb_index += 1
        candidates.append(
            {
                "id": f"autocodebenchmark_java_{acb_index:04d}",
                "source": "AutoCodeBenchmark",
                "source_record_id": item.get("question", "").splitlines()[0][:80],
                "source_path": "data/external_sources/java_cpp_self_contained/autocodebenchmark_java.jsonl",
                "code": code,
                "test_ref": test_ref,
                "question": item.get("question", ""),
                "difficulty": item.get("difficulty", ""),
                "effective_lines": effective_line_count(code),
                "complexity_score": len(code) / 80 + code.count("if (") + code.count("for (") + 2 * code.count("class "),
                "category_scores": scores,
                "category_eligibility": {
                    "simple_function": scores["simple_function"] >= 5,
                    "boundary": scores["boundary"] >= 4,
                    "complex_dependency": scores["complex_dependency"] >= 9,
                    "interface_mock": scores["interface_mock"] >= 10,
                },
            }
        )

    return candidates


def build_java_self_contained_dataset() -> None:
    ensure_dir(SELF_CONTAINED_OUTPUT_DIR)
    candidates = load_java_self_contained_candidates()
    selected, per_category = select_balanced_records(candidates, target_per_category=TARGET_PER_CATEGORY)
    counts = Counter(item["category"] for item in selected)
    missing = [category for category in CATEGORIES if counts[category] < TARGET_PER_CATEGORY]
    if missing:
        raise RuntimeError(
            "Not enough Java self-contained candidates: "
            + ", ".join(f"{category}={counts[category]}/{TARGET_PER_CATEGORY}" for category in missing)
        )

    samples = []
    for category in CATEGORIES:
        bucket_dir = SELF_CONTAINED_OUTPUT_DIR / category
        ensure_dir(bucket_dir)
        records = [item for item in selected if item["category"] == category]
        for index, record in enumerate(records):
            file_path = bucket_dir / f"{category}_{index:03d}.java"
            file_path.write_text(record["code"], encoding="utf-8")
            samples.append(
                {
                    "id": f"java_{category}_{index}",
                    "language": "java",
                    "category": category,
                    "source": record["source"],
                    "source_record_id": record["source_record_id"],
                    "source_path": record["source_path"],
                    "code_file": str(file_path.relative_to(BASE_DIR)).replace("\\", "/"),
                    "code": record["code"],
                    "test_ref": record["test_ref"],
                    "question": record["question"],
                    "difficulty": record["difficulty"],
                    "effective_lines": record["effective_lines"],
                    "category_scores": record["category_scores"],
                }
            )

    write_json(
        SELF_CONTAINED_OUTPUT_JSON,
        {
            "metadata": {
                "version": "Java_Self_Contained_v2",
                "language": "java",
                "description": "Balanced Java self-contained dataset built from HumanEval-X and AutoCodeBenchmark canonical solutions",
                "source_files": [
                    "data/external_sources/java_cpp_self_contained/humanevalx_java.jsonl",
                    "data/external_sources/java_cpp_self_contained/autocodebenchmark_java.jsonl",
                ],
                "selection_strategy": {
                    "per_category": TARGET_PER_CATEGORY,
                    "balanced_categories": True,
                    "require_self_contained_single_file": True,
                    "filters": [
                        "only java/javax imports",
                        "no package declaration",
                        "no third-party references",
                    ],
                },
            },
            "statistics": {
                "total": len(samples),
                "by_category": build_category_stats(samples),
                "by_source": dict(Counter(item["source"] for item in samples)),
                "candidate_pool": {category: len(per_category.get(category, [])) for category in CATEGORIES},
            },
            "samples": samples,
        },
    )


def main() -> None:
    build_java_module_dataset()
    build_java_self_contained_dataset()
    print("=== Java datasets ===")
    print(f"Self-contained output: {SELF_CONTAINED_OUTPUT_JSON}")
    print(f"Module-level output: {MODULE_OUTPUT_JSON}")


if __name__ == "__main__":
    main()
