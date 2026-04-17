import json
import re
from collections import Counter
from pathlib import Path

from dataset_builder_utils import CATEGORIES, build_category_stats, ensure_dir, select_balanced_records, write_json
from sync_external_self_contained_sources import OUTPUT_DIR as EXTERNAL_SOURCE_DIR
from sync_external_self_contained_sources import main as sync_external_sources


BASE_DIR = Path(__file__).resolve().parent
MODULE_SOURCE_FILE = BASE_DIR / "data" / "cpp_ut" / "cpp_ut_all.jsonl"
MODULE_OUTPUT_DIR = BASE_DIR / "data" / "cpp_code_files_module_level"
MODULE_OUTPUT_JSON = BASE_DIR / "data" / "cpp_ut_dataset" / "cpp_dataset_module_level.json"
SELF_CONTAINED_OUTPUT_DIR = BASE_DIR / "data" / "cpp_code_files_self_contained"
SELF_CONTAINED_OUTPUT_JSON = BASE_DIR / "data" / "cpp_ut_dataset" / "cpp_dataset_self_contained.json"
TARGET_PER_CATEGORY = 50


def effective_line_count(code: str) -> int:
    return sum(1 for line in code.splitlines() if line.strip())


def cpp_module_category_scores(code: str, test_ref: str) -> dict:
    low = f"{code}\n{test_ref}".lower()
    return {
        "interface_mock": (
            4 * sum(token in low for token in ["mock", "matcher", "listener", "callback", "spy", "stub", "fake"])
            + 2 * sum(token in low for token in ["virtual", "override", "template"])
            + low.count("matchandexplain")
        ),
        "complex_dependency": (
            3
            * sum(
                token in low
                for token in [
                    "protobuf",
                    "json",
                    "http",
                    "socket",
                    "stream",
                    "file",
                    "path",
                    "arena",
                    "memory",
                    "cache",
                    "absl",
                    "statusor",
                ]
            )
            + low.count("#include") / 4
        ),
        "boundary": (
            3
            * sum(
                token in low
                for token in [
                    "error",
                    "exception",
                    "invalid",
                    "null",
                    "empty",
                    "overflow",
                    "underflow",
                    "assert",
                    "range",
                    "failedprecondition",
                ]
            )
            + low.count("switch")
            + low.count("case ")
            + low.count("if (")
        ),
        "simple_function": (
            2 * sum(token in low for token in ["return ", "for (", "while (", "if ("])
            + (3 if all(token not in low for token in ["protobuf", "json", "http", "socket", "stream"]) else 0)
            - sum(token in low for token in ["mock", "matcher", "listener", "callback"])
        ),
    }


def build_cpp_module_dataset() -> None:
    records = []
    with MODULE_SOURCE_FILE.open("r", encoding="utf-8") as handle:
        for line in handle:
            if not line.strip():
                continue
            item = json.loads(line)
            code = item.get("Code", "")
            if not code.strip():
                continue
            effective_lines = effective_line_count(code)
            if not (15 <= effective_lines <= 260):
                continue
            test_ref = item.get("Unit Test - (Ground Truth)", "")
            scores = cpp_module_category_scores(code, test_ref)
            if max(scores.values()) <= 0:
                continue
            records.append(
                {
                    "id": str(item.get("ID")),
                    "language": "cpp",
                    "source": "CPP-UT-Bench",
                    "repository_name": item.get("Repository Name"),
                    "file_name": item.get("File Name"),
                    "source_path": item.get("File Path in Repository"),
                    "unit_test_path": item.get("File Path for Unit Test"),
                    "code": code,
                    "test_ref": test_ref,
                    "effective_lines": effective_lines,
                    "complexity_score": sum(scores.values()),
                    "category_scores": scores,
                    "category_eligibility": {category: scores[category] > 0 for category in CATEGORIES},
                }
            )

    selected, _ = select_balanced_records(records, target_per_category=TARGET_PER_CATEGORY)
    ensure_dir(MODULE_OUTPUT_DIR)
    samples = []
    for category in CATEGORIES:
        bucket = [item for item in selected if item["category"] == category]
        bucket_dir = MODULE_OUTPUT_DIR / category
        ensure_dir(bucket_dir)
        for index, record in enumerate(bucket):
            file_path = bucket_dir / f"{category}_{index:03d}.cpp"
            file_path.write_text(record["code"], encoding="utf-8")
            samples.append(
                {
                    "id": f"cpp_{category}_{index}",
                    "language": "cpp",
                    "category": category,
                    "source": record["source"],
                    "repository_name": record["repository_name"],
                    "file_name": record["file_name"],
                    "source_path": record["source_path"],
                    "unit_test_path": record["unit_test_path"],
                    "code_file": str(file_path.relative_to(BASE_DIR)).replace("\\", "/"),
                    "code": record["code"],
                    "test_ref": record["test_ref"],
                    "effective_lines": record["effective_lines"],
                    "category_scores": record["category_scores"],
                }
            )

    write_json(
        MODULE_OUTPUT_JSON,
        {
            "metadata": {
                "version": "CPP_Module_Level_v1",
                "language": "cpp",
                "source_file": str(MODULE_SOURCE_FILE.relative_to(BASE_DIR)).replace("\\", "/"),
                "description": "Balanced C++ project/module-level dataset rebuilt from cpp_ut_all.jsonl",
                "selection_strategy": {
                    "per_category": TARGET_PER_CATEGORY,
                    "balanced_categories": True,
                    "requires_project_context": True,
                    "require_self_contained_single_file": False,
                },
            },
            "statistics": {
                "total": len(samples),
                "by_category": build_category_stats(samples),
                "by_source": {"CPP-UT-Bench": len(samples)},
            },
            "samples": samples,
        },
    )


def load_jsonl(path: Path) -> list[dict]:
    with path.open("r", encoding="utf-8") as handle:
        return [json.loads(line) for line in handle if line.strip()]


def is_self_contained_cpp(code: str) -> tuple[bool, list[str]]:
    issues = []
    includes = re.findall(r"^\s*#include\s*([<\"].+[>\"])", code, re.M)
    for include in includes:
        if include.startswith('"'):
            issues.append(f"local_include:{include}")
            continue
        header = include[1:-1]
        if "/" in header and not header.startswith("bits/"):
            issues.append(f"non_std_include:{header}")
    if any(token in code for token in ["absl/", "gmock/", "gtest/", "boost/", "grpc/", "protobuf", "rapidjson", "nlohmann"]):
        issues.append("third_party_reference")
    return not issues, issues


def cpp_self_contained_scores(record: dict) -> dict[str, float]:
    code = record["code"]
    question = record.get("question", "")
    difficulty = record.get("difficulty", "")
    low = f"{question}\n{code}\n{record.get('test_ref', '')}".lower()
    class_count = len(re.findall(r"\bclass\s+\w+", code))
    method_tokens = sum(token in low for token in ["vector<", "map<", "unordered_map", "string", "algorithm", "stack", "queue"])
    source_bias = 4 if record["source"] == "HumanEval-X" else 0
    oop_bias = 5 if record["source"] == "AutoCodeBenchmark" else 0

    return {
        "simple_function": (
            3 * sum(token in low for token in ["return ", "for (", "while (", "vector<", "string", "sort"])
            + source_bias
            + (2 if difficulty == "easy" else 0)
            - class_count
        ),
        "boundary": (
            4 * sum(token in low for token in ["throw", "exception", "invalid", "nullptr", "empty", "range", "overflow"])
            + low.count("if (")
            + source_bias
            + (2 if difficulty == "medium" else 0)
        ),
        "complex_dependency": (
            4 * method_tokens
            + 4 * sum(token in low for token in ["matrix", "graph", "tree", "priority_queue", "unordered_map", "set<"])
            + oop_bias
            + (3 if difficulty == "hard" else 0)
        ),
        "interface_mock": (
            5 * sum(token in low for token in ["class ", "manager", "system", "analyzer", "account", "user", "role"])
            + 2 * class_count
            + oop_bias
            + (3 if difficulty == "hard" else 0)
        ),
    }


def load_cpp_self_contained_candidates() -> list[dict]:
    sync_external_sources()
    candidates = []

    for item in load_jsonl(EXTERNAL_SOURCE_DIR / "humanevalx_cpp.jsonl"):
        code = (item["prompt"] + item["canonical_solution"]).strip() + "\n"
        ok, issues = is_self_contained_cpp(code)
        if not ok:
            continue
        scores = cpp_self_contained_scores(
            {
                "source": "HumanEval-X",
                "code": code,
                "test_ref": item.get("test", ""),
                "question": item.get("prompt", ""),
                "difficulty": "benchmark",
            }
        )
        candidates.append(
            {
                "id": f"humanevalx_{item['task_id'].replace('/', '_').lower()}",
                "source": "HumanEval-X",
                "source_record_id": item["task_id"],
                "source_path": "data/external_sources/java_cpp_self_contained/humanevalx_cpp.jsonl",
                "code": code,
                "test_ref": item.get("test", ""),
                "question": item.get("prompt", ""),
                "difficulty": "benchmark",
                "effective_lines": effective_line_count(code),
                "complexity_score": len(code) / 70 + code.count("if (") + code.count("for (") + 2 * code.count("class "),
                "category_scores": scores,
                "category_eligibility": {
                    "simple_function": scores["simple_function"] >= 6,
                    "boundary": scores["boundary"] >= 5,
                    "complex_dependency": scores["complex_dependency"] >= 4,
                    "interface_mock": scores["interface_mock"] >= 4,
                },
            }
        )

    acb_index = 0
    for item in load_jsonl(EXTERNAL_SOURCE_DIR / "autocodebenchmark_cpp.jsonl"):
        code = item["canonical_solution"].strip()
        if not code:
            continue
        code = code + "\n"
        ok, issues = is_self_contained_cpp(code)
        if not ok:
            continue
        test_ref = item.get("full_test_func") or item.get("demo_test_func") or ""
        scores = cpp_self_contained_scores(
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
                "id": f"autocodebenchmark_cpp_{acb_index:04d}",
                "source": "AutoCodeBenchmark",
                "source_record_id": item.get("question", "").splitlines()[0][:80],
                "source_path": "data/external_sources/java_cpp_self_contained/autocodebenchmark_cpp.jsonl",
                "code": code,
                "test_ref": test_ref,
                "question": item.get("question", ""),
                "difficulty": item.get("difficulty", ""),
                "effective_lines": effective_line_count(code),
                "complexity_score": len(code) / 70 + code.count("if (") + code.count("for (") + 2 * code.count("class "),
                "category_scores": scores,
                "category_eligibility": {
                    "simple_function": scores["simple_function"] >= 5,
                    "boundary": scores["boundary"] >= 4,
                    "complex_dependency": scores["complex_dependency"] >= 10,
                    "interface_mock": scores["interface_mock"] >= 10,
                },
            }
        )

    return candidates


def build_cpp_self_contained_dataset() -> None:
    ensure_dir(SELF_CONTAINED_OUTPUT_DIR)
    candidates = load_cpp_self_contained_candidates()
    selected, per_category = select_balanced_records(candidates, target_per_category=TARGET_PER_CATEGORY)
    counts = Counter(item["category"] for item in selected)
    missing = [category for category in CATEGORIES if counts[category] < TARGET_PER_CATEGORY]
    if missing:
        raise RuntimeError(
            "Not enough C++ self-contained candidates: "
            + ", ".join(f"{category}={counts[category]}/{TARGET_PER_CATEGORY}" for category in missing)
        )

    samples = []
    for category in CATEGORIES:
        bucket_dir = SELF_CONTAINED_OUTPUT_DIR / category
        ensure_dir(bucket_dir)
        records = [item for item in selected if item["category"] == category]
        for index, record in enumerate(records):
            file_path = bucket_dir / f"{category}_{index:03d}.cpp"
            file_path.write_text(record["code"], encoding="utf-8")
            samples.append(
                {
                    "id": f"cpp_{category}_{index}",
                    "language": "cpp",
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
                "version": "CPP_Self_Contained_v2",
                "language": "cpp",
                "description": "Balanced C++ self-contained dataset built from HumanEval-X and AutoCodeBenchmark canonical solutions",
                "source_files": [
                    "data/external_sources/java_cpp_self_contained/humanevalx_cpp.jsonl",
                    "data/external_sources/java_cpp_self_contained/autocodebenchmark_cpp.jsonl",
                ],
                "selection_strategy": {
                    "per_category": TARGET_PER_CATEGORY,
                    "balanced_categories": True,
                    "require_self_contained_single_file": True,
                    "filters": [
                        "no local includes",
                        "only standard library headers",
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
    build_cpp_module_dataset()
    build_cpp_self_contained_dataset()
    print("=== C++ datasets ===")
    print(f"Self-contained output: {SELF_CONTAINED_OUTPUT_JSON}")
    print(f"Module-level output: {MODULE_OUTPUT_JSON}")


if __name__ == "__main__":
    main()
