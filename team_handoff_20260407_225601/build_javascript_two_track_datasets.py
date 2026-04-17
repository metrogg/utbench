import re
from collections import Counter, defaultdict
from pathlib import Path

from dataset_builder_utils import CATEGORIES, ensure_dir, write_json


BASE_DIR = Path(__file__).resolve().parent
SOURCE_FILE = BASE_DIR / "js_sources" / "lodash.js"

SELF_CONTAINED_DIR = BASE_DIR / "data" / "javascript_code_files_self_contained"
SELF_CONTAINED_JSON = BASE_DIR / "data" / "javascript_ut_dataset" / "javascript_dataset_self_contained.json"
MODULE_LEVEL_DIR = BASE_DIR / "data" / "javascript_code_files_module_level"
MODULE_LEVEL_JSON = BASE_DIR / "data" / "javascript_ut_dataset" / "javascript_dataset_module_level.json"

TARGET_PER_CATEGORY = 50
SELECTION_ORDER = ["interface_mock", "complex_dependency", "boundary", "simple_function"]

FUNCTION_RE = re.compile(r"function\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\([^)]*\)\s*\{")
CALL_RE = re.compile(r"\b([A-Za-z_$][A-Za-z0-9_$]*)\s*\(")
KEYWORDS = {"if", "for", "while", "switch", "catch", "function", "return"}
ALLOWED_EXTERNAL_CALLS = {
    "Array",
    "Boolean",
    "Date",
    "Error",
    "Math",
    "Number",
    "Object",
    "RegExp",
    "String",
    "SyntaxError",
    "TypeError",
    "charAt",
    "clearInterval",
    "clearTimeout",
    "hasOwnProperty",
    "isFinite",
    "isNaN",
    "match",
    "parseFloat",
    "parseInt",
    "push",
    "replace",
    "setInterval",
    "setTimeout",
    "slice",
    "sort",
    "test",
    "toString",
}


def extract_functions(source_text: str):
    functions = {}
    order = []
    for match in FUNCTION_RE.finditer(source_text):
        name = match.group(1)
        start = match.start()
        body_start = match.end() - 1
        depth = 0
        pos = body_start
        while pos < len(source_text):
            char = source_text[pos]
            if char == "{":
                depth += 1
            elif char == "}":
                depth -= 1
                if depth == 0:
                    end = pos + 1
                    functions[name] = source_text[start:end]
                    order.append(name)
                    break
            pos += 1
    return functions, order


def category_scores(name: str, code: str) -> dict[str, float]:
    low = code.lower()
    dependency_markers = [
        "cache",
        "wrapper",
        "clone",
        "merge",
        "assign",
        "path",
        "keys",
        "values",
        "stack",
        "map",
        "set",
        "hash",
        "flatten",
    ]
    interface_markers = ["iteratee", "callback", "comparator", "customizer", "predicate", "func", "wrapper"]
    boundary_markers = ["null", "undefined", "nan", "trim", "unicode", "ascii", "empty", "length", "index"]
    return {
        "simple_function": (
            2 * sum(token in low for token in ["return ", "while ", "for ", "if "])
            + (3 if "cache" not in low and "wrapper" not in low else 0)
            - sum(token in low for token in interface_markers)
        ),
        "boundary": (
            3 * sum(token in low for token in boundary_markers)
            + low.count("if (")
            + low.count("=== undefined")
            + low.count("== null")
        ),
        "interface_mock": (
            4 * sum(token in low for token in ["iteratee", "callback", "comparator", "customizer", "predicate"])
            + 2 * sum(token in low for token in ["func", "wrapper"])
            + low.count("call(")
            + low.count("apply(")
        ),
        "complex_dependency": (
            3 * sum(token in low for token in dependency_markers)
            + low.count("new ")
            + low.count("get(")
            + low.count("set(")
        ),
    }


def category_eligibility(scores: dict[str, float]) -> dict[str, bool]:
    return {category: scores[category] > 0 for category in CATEGORIES}


def analyze_function_calls(function_name: str, code: str, function_names: set[str]):
    params_match = re.search(rf"function\s+{re.escape(function_name)}\s*\(([^)]*)\)", code)
    params = {part.strip() for part in params_match.group(1).split(",") if part.strip()} if params_match else set()
    calls = {name for name in CALL_RE.findall(code) if name not in KEYWORDS and name != function_name}
    local_dependencies = sorted(name for name in calls if name in function_names and name not in params)
    external_calls = sorted(
        name
        for name in calls
        if name not in function_names and name not in params and name not in ALLOWED_EXTERNAL_CALLS
    )
    return local_dependencies, external_calls


def render_closure_bundle(root_name: str, functions: dict[str, str], order: list[str]):
    visited = set()
    external_issues = set()

    def visit(name: str):
        if name in visited:
            return
        visited.add(name)
        local_deps, external_calls = analyze_function_calls(name, functions[name], set(functions))
        external_issues.update(external_calls)
        for dep in local_deps:
            visit(dep)

    visit(root_name)
    bundle = "\n\n".join(functions[name] for name in order if name in visited)
    return bundle, sorted(visited), sorted(external_issues)


def build_module_level_records():
    source_text = SOURCE_FILE.read_text(encoding="utf-8")
    functions, order = extract_functions(source_text)
    candidates = []
    for name in order:
        code = functions[name]
        effective_lines = sum(1 for line in code.splitlines() if line.strip())
        if not 8 <= effective_lines <= 140:
            continue
        scores = category_scores(name, code)
        eligibility = category_eligibility(scores)
        if not any(eligibility.values()):
            continue
        candidates.append(
            {
                "id": name,
                "language": "javascript",
                "source": "lodash_module",
                "project": "lodash",
                "source_path": "js_sources/lodash.js",
                "function_name": name,
                "code": code,
                "effective_lines": effective_lines,
                "complexity_score": sum(scores.values()),
                "category_scores": scores,
                "category_eligibility": eligibility,
            }
        )

    by_category = defaultdict(list)
    for candidate in candidates:
        for category in CATEGORIES:
            if candidate["category_eligibility"][category]:
                by_category[category].append(candidate)
    for category in CATEGORIES:
        by_category[category].sort(
            key=lambda item: (-item["category_scores"][category], -item["effective_lines"], item["function_name"])
        )

    selected = []
    used_names = set()
    for category in SELECTION_ORDER:
        picked = 0
        for candidate in by_category[category]:
            if picked >= TARGET_PER_CATEGORY:
                break
            if candidate["id"] in used_names:
                continue
            selected.append({**candidate, "category": category})
            used_names.add(candidate["id"])
            picked += 1
    return selected, by_category


def build_self_contained_records():
    source_text = SOURCE_FILE.read_text(encoding="utf-8")
    functions, order = extract_functions(source_text)
    records = []
    for name in order:
        bundle, bundled_functions, external_issues = render_closure_bundle(name, functions, order)
        if external_issues:
            continue
        effective_lines = sum(1 for line in bundle.splitlines() if line.strip())
        if not 8 <= effective_lines <= 220:
            continue
        scores = category_scores(name, bundle)
        eligibility = category_eligibility(scores)
        if not any(eligibility.values()):
            continue
        records.append(
            {
                "id": name,
                "language": "javascript",
                "source": "lodash_self_contained",
                "project": "lodash",
                "source_path": "js_sources/lodash.js",
                "function_name": name,
                "bundled_functions": bundled_functions,
                "code": bundle,
                "effective_lines": effective_lines,
                "complexity_score": sum(scores.values()),
                "category_scores": scores,
                "category_eligibility": eligibility,
                "category": max(scores, key=scores.get),
                "self_contained_check": {
                    "is_self_contained": True,
                    "bundled_functions": bundled_functions,
                    "external_issues": [],
                },
            }
        )
    return records


def write_dataset(records, *, output_dir: Path, output_json: Path, version: str, description: str, selection_strategy: dict):
    ensure_dir(output_dir)
    samples = []
    by_category_stats = {}
    for category in CATEGORIES:
        bucket = [item for item in records if item["category"] == category]
        bucket_dir = output_dir / category
        ensure_dir(bucket_dir)
        lengths = [len(item["code"]) for item in bucket]
        by_category_stats[category] = {
            "count": len(bucket),
            "avg_effective_lines": round(sum(item["effective_lines"] for item in bucket) / len(bucket), 2)
            if bucket
            else 0,
            "min_length": min(lengths) if lengths else 0,
            "max_length": max(lengths) if lengths else 0,
        }
        for index, record in enumerate(bucket):
            file_path = bucket_dir / f"{category}_{index:03d}.js"
            file_path.write_text(record["code"], encoding="utf-8")
            sample = {
                "id": f"javascript_{category}_{index}",
                "language": "javascript",
                "category": category,
                "source": record["source"],
                "project": record["project"],
                "source_path": record["source_path"],
                "function_name": record["function_name"],
                "code_file": str(file_path.relative_to(BASE_DIR)).replace("\\", "/"),
                "code": record["code"],
                "effective_lines": record["effective_lines"],
                "category_scores": record["category_scores"],
            }
            if "bundled_functions" in record:
                sample["bundled_functions"] = record["bundled_functions"]
            if "self_contained_check" in record:
                sample["self_contained_check"] = record["self_contained_check"]
            samples.append(sample)

    payload = {
        "metadata": {
            "version": version,
            "language": "javascript",
            "source_file": str(SOURCE_FILE.relative_to(BASE_DIR)).replace("\\", "/"),
            "description": description,
            "selection_strategy": selection_strategy,
        },
        "statistics": {
            "total": len(samples),
            "by_category": by_category_stats,
            "by_source": dict(Counter(item["source"] for item in samples)),
            "projects": dict(Counter(item["project"] for item in samples)),
        },
        "samples": samples,
    }
    write_json(output_json, payload)


def main():
    module_records, module_pools = build_module_level_records()
    self_contained_records = build_self_contained_records()

    write_dataset(
        module_records,
        output_dir=MODULE_LEVEL_DIR,
        output_json=MODULE_LEVEL_JSON,
        version="JavaScript_Module_Level_v1",
        description="Balanced JavaScript dataset extracted from lodash.js as project-linked module slices",
        selection_strategy={
            "per_category": TARGET_PER_CATEGORY,
            "balanced_categories": True,
            "requires_project_context": True,
        },
    )
    write_dataset(
        self_contained_records,
        output_dir=SELF_CONTAINED_DIR,
        output_json=SELF_CONTAINED_JSON,
        version="JavaScript_Self_Contained_v1",
        description="Self-contained JavaScript dataset built by inlining lodash helper-function closures",
        selection_strategy={
            "balanced_categories": False,
            "require_self_contained_single_file": True,
            "allow_project_local_helper_inlining": True,
        },
    )

    print("=== JavaScript datasets ===")
    print(f"Module-level records: {len(module_records)}")
    print(f"Self-contained records: {len(self_contained_records)}")
    for category in CATEGORIES:
        print(f"{category}: module pool {len(module_pools[category])}")
    print(f"Module-level output: {MODULE_LEVEL_JSON}")
    print(f"Self-contained output: {SELF_CONTAINED_JSON}")


if __name__ == "__main__":
    main()
