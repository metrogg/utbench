import json
import re
from collections import Counter
from pathlib import Path

from dataset_builder_utils import CATEGORIES, ensure_dir, select_balanced_records, write_json


BASE_DIR = Path(__file__).resolve().parent
SOURCE_FILE = BASE_DIR / "data" / "go" / "go_all.jsonl"

SELF_CONTAINED_DIR = BASE_DIR / "data" / "go_code_files_self_contained"
SELF_CONTAINED_JSON = BASE_DIR / "data" / "go_ut_dataset" / "go_dataset_self_contained.json"
MODULE_LEVEL_DIR = BASE_DIR / "data" / "go_code_files_module_level"
MODULE_LEVEL_JSON = BASE_DIR / "data" / "go_ut_dataset" / "go_dataset_module_level.json"

TARGET_PER_CATEGORY = 50
SELECTION_ORDER = ["interface_mock", "complex_dependency", "boundary", "simple_function"]

GO_STDLIB_IMPORTS = {
    "bufio",
    "bytes",
    "context",
    "crypto/tls",
    "crypto/x509",
    "encoding/base64",
    "encoding/binary",
    "encoding/csv",
    "encoding/hex",
    "encoding/json",
    "encoding/pem",
    "errors",
    "fmt",
    "io",
    "io/ioutil",
    "math",
    "math/rand",
    "net/http",
    "net/url",
    "os",
    "path",
    "path/filepath",
    "regexp",
    "sort",
    "strconv",
    "strings",
    "sync",
    "sync/atomic",
    "time",
    "unicode",
    "unicode/utf8",
}
GO_QUALIFIER_TO_IMPORT = {
    "base64": "encoding/base64",
    "binary": "encoding/binary",
    "bufio": "bufio",
    "bytes": "bytes",
    "context": "context",
    "csv": "encoding/csv",
    "errors": "errors",
    "fmt": "fmt",
    "hex": "encoding/hex",
    "http": "net/http",
    "io": "io",
    "ioutil": "io/ioutil",
    "json": "encoding/json",
    "math": "math",
    "os": "os",
    "path": "path",
    "filepath": "path/filepath",
    "pem": "encoding/pem",
    "rand": "math/rand",
    "regexp": "regexp",
    "sort": "sort",
    "strconv": "strconv",
    "strings": "strings",
    "sync": "sync",
    "atomic": "sync/atomic",
    "time": "time",
    "tls": "crypto/tls",
    "unicode": "unicode",
    "utf8": "unicode/utf8",
    "url": "net/url",
    "x509": "crypto/x509",
}
GO_ALLOWED_TYPE_IDENTIFIERS = {
    "bool",
    "byte",
    "complex128",
    "complex64",
    "error",
    "float32",
    "float64",
    "int",
    "int16",
    "int32",
    "int64",
    "int8",
    "rune",
    "string",
    "uint",
    "uint16",
    "uint32",
    "uint64",
    "uint8",
    "uintptr",
    "Buffer",
    "Context",
    "Duration",
    "File",
    "IP",
    "IPNet",
    "Reader",
    "ReadCloser",
    "Regexp",
    "Scanner",
    "Stringer",
    "Time",
    "Timer",
    "Ticker",
    "URL",
    "Writer",
    "WriteCloser",
    "CertPool",
    "Certificate",
}

FUNC_SIGNATURE_RE = re.compile(r"^func\s*(\([^)]*\)\s*)?([A-Za-z_][A-Za-z0-9_]*)\s*\((.*?)\)\s*(.*?)\s*\{", re.S)
QUALIFIER_RE = re.compile(r"\b([a-z_][A-Za-z0-9_]*)\.")


def read_go_records():
    with SOURCE_FILE.open("r", encoding="utf-8") as handle:
        for line in handle:
            if not line.strip():
                continue
            yield json.loads(line)


def identifiers_from_go_type(type_text: str) -> list[str]:
    return [
        token
        for token in re.findall(r"[A-Za-z_][A-Za-z0-9_\.]*", type_text)
        if token not in {"map", "chan", "func", "interface", "struct"}
    ]


def analyze_go_self_contained(function_code: str) -> tuple[bool, list[str], list[str]]:
    issues: list[str] = []
    imports: set[str] = set()

    signature_match = FUNC_SIGNATURE_RE.search(function_code)
    if not signature_match:
        return False, [], ["bad_signature"]

    receiver, _, params, results = signature_match.groups()
    if receiver:
        issues.append("method_receiver")

    qualifiers = set(QUALIFIER_RE.findall(function_code))
    if "testing" in qualifiers or "testing." in params or "testing." in results:
        issues.append("testing_dependency")

    for qualifier in qualifiers:
        import_path = GO_QUALIFIER_TO_IMPORT.get(qualifier)
        if not import_path:
            issues.append(f"non_stdlib_qualifier:{qualifier}")
        else:
            imports.add(import_path)

    signature_text = f"{params} {results}"
    for identifier in identifiers_from_go_type(signature_text):
        if "." in identifier:
            pkg, _ = identifier.split(".", 1)
            if pkg not in GO_QUALIFIER_TO_IMPORT:
                issues.append(f"non_stdlib_signature_type:{identifier}")
        elif identifier not in GO_ALLOWED_TYPE_IDENTIFIERS:
            issues.append(f"custom_signature_type:{identifier}")

    return not issues, sorted(imports), issues


def render_go_self_contained(function_code: str, imports: list[str]) -> str:
    parts = ["package main", ""]
    if imports:
        parts.append("import (")
        parts.extend(f'\t"{import_path}"' for import_path in imports)
        parts.append(")")
        parts.append("")
    parts.append(function_code.rstrip() + "\n")
    return "\n".join(parts)


def go_category_scores(code: str, path_text: str) -> dict[str, float]:
    low = f"{path_text}\n{code}".lower()
    dependency_hits = sum(
        token in low
        for token in [
            "os.",
            "io.",
            "ioutil.",
            "filepath.",
            "path.",
            "http.",
            "json.",
            "xml.",
            "tls.",
            "x509.",
            "pem.",
            "url.",
            "context.",
            "time.",
        ]
    )
    boundary_hits = sum(
        token in low
        for token in [
            "if err != nil",
            "return nil, err",
            "panic(",
            "invalid",
            "error",
            " nil",
            "parse",
            "validate",
            "switch ",
            "case ",
        ]
    )
    interface_hits = sum(
        token in low
        for token in [
            "interface{",
            " interface ",
            "handler",
            "callback",
            "middleware",
            "client",
            "server",
            "reader",
            "writer",
        ]
    )
    simple_hits = low.count("return ") + low.count("for ") + low.count("if ") + low.count("len(")
    return {
        "simple_function": simple_hits * 2 + 3 - dependency_hits - interface_hits,
        "boundary": boundary_hits * 3 + low.count("if ") + low.count("switch "),
        "complex_dependency": dependency_hits * 4 + low.count("new") + low.count("make("),
        "interface_mock": interface_hits * 4 + low.count("func(") * 2,
    }


def go_category_eligibility(scores: dict[str, float]) -> dict[str, bool]:
    return {
        "simple_function": scores["simple_function"] > 0,
        "boundary": scores["boundary"] > 0,
        "complex_dependency": scores["complex_dependency"] > 0,
        "interface_mock": scores["interface_mock"] > 0,
    }


def build_candidates():
    module_candidates = []
    self_contained_candidates = []
    for item in read_go_records():
        code = item["whole_func_string"].strip() + "\n"
        path_text = item["func_path_in_repository"]
        scores = go_category_scores(code, path_text)
        eligibility = go_category_eligibility(scores)
        if not any(eligibility.values()):
            continue

        base_record = {
            "id": item["func_code_url"],
            "language": "go",
            "source": "project_function",
            "repository_name": item["repository_name"],
            "source_path": path_text,
            "function_name": item["func_name"],
            "docstring": item.get("func_documentation_string", ""),
            "code": code,
            "effective_lines": sum(1 for line in code.splitlines() if line.strip()),
            "complexity_score": sum(scores.values()),
            "category_scores": scores,
            "category_eligibility": eligibility,
        }
        module_candidates.append(base_record)

        ok, imports, issues = analyze_go_self_contained(code)
        if not ok:
            continue
        rendered = render_go_self_contained(code, imports)
        self_contained_candidates.append(
            {
                **base_record,
                "code": rendered,
                "effective_lines": sum(1 for line in rendered.splitlines() if line.strip()),
                "self_contained_check": {
                    "is_self_contained": True,
                    "imports": imports,
                    "issues": [],
                },
            }
        )
    return module_candidates, self_contained_candidates


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
            file_path = bucket_dir / f"{category}_{index:03d}.go"
            file_path.write_text(record["code"], encoding="utf-8")
            sample = {
                "id": f"go_{category}_{index}",
                "language": "go",
                "category": category,
                "source": record["source"],
                "repository_name": record["repository_name"],
                "source_path": record["source_path"],
                "function_name": record["function_name"],
                "code_file": str(file_path.relative_to(BASE_DIR)).replace("\\", "/"),
                "code": record["code"],
                "effective_lines": record["effective_lines"],
                "category_scores": record["category_scores"],
            }
            if "self_contained_check" in record:
                sample["self_contained_check"] = record["self_contained_check"]
            samples.append(sample)

    payload = {
        "metadata": {
            "version": version,
            "language": "go",
            "source_file": str(SOURCE_FILE.relative_to(BASE_DIR)).replace("\\", "/"),
            "description": description,
            "selection_strategy": selection_strategy,
        },
        "statistics": {
            "total": len(samples),
            "by_category": by_category_stats,
            "by_source": dict(Counter(item["source"] for item in samples)),
            "repositories": dict(Counter(item["repository_name"] for item in samples)),
        },
        "samples": samples,
    }
    write_json(output_json, payload)


def main():
    module_candidates, self_contained_candidates = build_candidates()

    selected_module, pools_module = select_balanced_records(
        module_candidates,
        categories=CATEGORIES,
        target_per_category=TARGET_PER_CATEGORY,
        selection_order=SELECTION_ORDER,
        unique_key="id",
    )
    selected_self, pools_self = select_balanced_records(
        self_contained_candidates,
        categories=CATEGORIES,
        target_per_category=TARGET_PER_CATEGORY,
        selection_order=SELECTION_ORDER,
        unique_key="id",
    )

    write_dataset(
        selected_module,
        output_dir=MODULE_LEVEL_DIR,
        output_json=MODULE_LEVEL_JSON,
        version="GO_Module_Level_v1",
        description="Balanced Go dataset built from real project functions with repository metadata",
        selection_strategy={
            "per_category": TARGET_PER_CATEGORY,
            "granularity": "project_function",
            "requires_project_context": True,
            "balanced_categories": True,
        },
    )
    write_dataset(
        selected_self,
        output_dir=SELF_CONTAINED_DIR,
        output_json=SELF_CONTAINED_JSON,
        version="GO_Self_Contained_v1",
        description="Balanced Go dataset built from stdlib-only real functions rendered as standalone files",
        selection_strategy={
            "per_category": TARGET_PER_CATEGORY,
            "granularity": "standalone_function",
            "require_self_contained_single_file": True,
            "allow_stdlib_imports_only": True,
            "balanced_categories": True,
        },
    )

    print("=== Go datasets ===")
    print(f"Module-level candidates: {len(module_candidates)}")
    print(f"Self-contained candidates: {len(self_contained_candidates)}")
    for category in CATEGORIES:
        print(
            f"{category}: module {len(pools_module[category])} pool, self-contained {len(pools_self[category])} pool"
        )
    print(f"Module-level output: {MODULE_LEVEL_JSON}")
    print(f"Self-contained output: {SELF_CONTAINED_JSON}")


if __name__ == "__main__":
    main()
