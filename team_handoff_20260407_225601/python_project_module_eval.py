import argparse
import builtins
import json
import os
import shutil
import stat
import subprocess
import sys
import xml.etree.ElementTree as ET
from datetime import datetime
from pathlib import Path


BASE_DIR = Path(__file__).resolve().parent
DEFAULT_DATASET = BASE_DIR / "data" / "python_ut_dataset" / "python_dataset_module_level.json"
DEFAULT_RUN_ROOT = BASE_DIR / "eval" / "project_level_runs"


def load_dataset(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def find_sample(dataset: dict, sample_id: str | None, sample_index: int | None) -> dict:
    samples = dataset["samples"]
    if sample_id:
        for sample in samples:
            if sample["id"] == sample_id:
                return sample
        raise KeyError(f"Sample id not found: {sample_id}")
    if sample_index is not None:
        return samples[sample_index]
    return samples[0]


def extract_code_block(text: str) -> str:
    stripped = text.strip()
    if not stripped.startswith("```"):
        return stripped
    lines = stripped.splitlines()
    if lines and lines[0].startswith("```"):
        lines = lines[1:]
    if lines and lines[-1].strip() == "```":
        lines = lines[:-1]
    return "\n".join(lines).strip() + "\n"


def resolve_package_copy(source_path: Path, package: str) -> tuple[Path, Path, str]:
    package_parts = [part for part in package.split(".") if part]
    parts = list(source_path.parts)
    package_index = -1

    if package_parts:
        for index in range(len(parts) - len(package_parts)):
            if parts[index : index + len(package_parts)] == package_parts:
                package_index = index
                break

    if package_index >= 0:
        package_root = Path(*parts[:package_index])
        package_dir = package_root.joinpath(*package_parts)
        module_rel = source_path.relative_to(package_root)
        module_import = ".".join(module_rel.with_suffix("").parts)
        if module_import.endswith(".__init__"):
            module_import = module_import[: -len(".__init__")]
        return package_dir, module_rel, module_import

    package_dir = source_path.parent
    module_rel = Path(source_path.name)
    module_import = source_path.stem
    return package_dir, module_rel, module_import


def restore_workspace(sample: dict, run_dir: Path) -> dict:
    source_path = BASE_DIR / sample["source_path"]
    if not source_path.exists():
        raise FileNotFoundError(f"Source file not found: {source_path}")

    package_dir, module_rel, module_import = resolve_package_copy(source_path, sample.get("package", ""))
    workspace = run_dir / "workspace"
    workspace.mkdir(parents=True, exist_ok=True)

    destination = workspace / package_dir.name
    shutil.copytree(
        package_dir,
        destination,
        ignore=shutil.ignore_patterns("__pycache__", "*.pyc", "*.pyo"),
    )
    for copied_path in destination.rglob("*"):
        try:
            os.chmod(copied_path, stat.S_IWRITE | stat.S_IREAD)
        except OSError:
            pass
    target_file = workspace / module_rel
    tests_dir = workspace / "tests"
    tests_dir.mkdir(parents=True, exist_ok=True)

    return {
        "workspace": workspace,
        "target_file": target_file,
        "tests_dir": tests_dir,
        "module_import": module_import,
        "package_name": package_dir.name,
        "source_path": source_path,
    }


def run_syntax_compile(files: list[Path]) -> tuple[bool, str]:
    errors = []
    for path in files:
        try:
            source = path.read_text(encoding="utf-8")
            builtins.compile(source, str(path), "exec")
        except Exception as exc:
            errors.append(f"{path}: {exc}")
    return len(errors) == 0, "\n".join(errors)


def run_pytest(
    python_exe: Path,
    workspace: Path,
    tests_file: Path,
    package_name: str,
    junit_xml: Path,
    coverage_data: Path,
) -> subprocess.CompletedProcess:
    env = os.environ.copy()
    env["PYTHONPATH"] = str(workspace)
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    env["COVERAGE_FILE"] = str(coverage_data)
    command = [
        str(python_exe),
        "-m",
        "coverage",
        "run",
        "--branch",
        "--source",
        package_name,
        "-m",
        "pytest",
        str(tests_file),
        "-q",
        "-p",
        "no:cacheprovider",
        f"--junitxml={junit_xml}",
    ]
    return subprocess.run(command, cwd=workspace, env=env, capture_output=True, text=True)


def run_coverage_json(
    python_exe: Path,
    workspace: Path,
    coverage_data: Path,
    coverage_json: Path,
) -> subprocess.CompletedProcess:
    env = os.environ.copy()
    env["PYTHONPATH"] = str(workspace)
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    env["COVERAGE_FILE"] = str(coverage_data)
    command = [
        str(python_exe),
        "-m",
        "coverage",
        "json",
        "-o",
        str(coverage_json),
    ]
    return subprocess.run(command, cwd=workspace, env=env, capture_output=True, text=True)


def parse_junit_xml(path: Path) -> dict:
    if not path.exists():
        return {
            "tests": 0,
            "failures": 0,
            "errors": 0,
            "skipped": 0,
            "passed": 0,
        }

    root = ET.fromstring(path.read_text(encoding="utf-8"))
    suite = root if root.tag == "testsuite" else root.find("testsuite")
    if suite is None:
        return {
            "tests": 0,
            "failures": 0,
            "errors": 0,
            "skipped": 0,
            "passed": 0,
        }

    tests = int(suite.attrib.get("tests", 0))
    failures = int(suite.attrib.get("failures", 0))
    errors = int(suite.attrib.get("errors", 0))
    skipped = int(suite.attrib.get("skipped", 0))
    passed = max(tests - failures - errors - skipped, 0)
    return {
        "tests": tests,
        "failures": failures,
        "errors": errors,
        "skipped": skipped,
        "passed": passed,
    }


def parse_coverage(path: Path, target_file: Path) -> dict:
    if not path.exists():
        return {
            "measured": False,
            "target_file": str(target_file),
        }

    payload = json.loads(path.read_text(encoding="utf-8"))
    normalized_target = str(target_file).replace("\\", "/")
    match = None

    for file_path, stats in payload.get("files", {}).items():
        if file_path.replace("\\", "/").endswith(normalized_target):
            match = (file_path, stats)
            break
        if file_path.replace("\\", "/").endswith(target_file.name):
            match = (file_path, stats)

    if match is None:
        return {
            "measured": False,
            "target_file": str(target_file),
            "available_files": sorted(payload.get("files", {}).keys())[:20],
        }

    file_path, stats = match
    summary = stats.get("summary", {})
    return {
        "measured": True,
        "target_file": file_path,
        "covered_lines": summary.get("covered_lines", 0),
        "num_statements": summary.get("num_statements", 0),
        "missing_lines": summary.get("missing_lines", 0),
        "excluded_lines": summary.get("excluded_lines", 0),
        "num_branches": summary.get("num_branches", 0),
        "covered_branches": summary.get("covered_branches", 0),
        "missing_branches": summary.get("missing_branches", 0),
        "percent_covered": summary.get("percent_covered", 0.0),
        "percent_covered_display": summary.get("percent_covered_display", "0"),
        "executed_lines_sample": stats.get("executed_lines", [])[:30],
        "missing_lines_sample": stats.get("missing_lines", [])[:30],
    }


def save_report(run_dir: Path, payload: dict) -> Path:
    run_dir.mkdir(parents=True, exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    report_path = run_dir / f"report_{timestamp}.json"
    report_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    return report_path


def main() -> int:
    parser = argparse.ArgumentParser(description="Restore and evaluate a Python module-level sample in project context.")
    parser.add_argument("--dataset", type=Path, default=DEFAULT_DATASET)
    parser.add_argument("--sample-id", type=str, default=None)
    parser.add_argument("--sample-index", type=int, default=None)
    parser.add_argument("--test-file", type=Path, required=True)
    parser.add_argument("--model-name", type=str, default="manual_model")
    parser.add_argument("--python-exe", type=Path, default=BASE_DIR / "venv" / "Scripts" / "python.exe")
    parser.add_argument("--run-root", type=Path, default=DEFAULT_RUN_ROOT)
    args = parser.parse_args()

    dataset = load_dataset(args.dataset)
    sample = find_sample(dataset, args.sample_id, args.sample_index)
    run_stamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    run_dir = args.run_root / sample["id"] / run_stamp
    restored = restore_workspace(sample, run_dir)

    generated_text = extract_code_block(args.test_file.read_text(encoding="utf-8"))
    tests_file = restored["tests_dir"] / f"test_generated_{Path(sample['source_path']).stem}.py"
    tests_file.write_text(generated_text, encoding="utf-8")

    compile_ok, compile_output = run_syntax_compile(
        [restored["target_file"], tests_file],
    )

    junit_xml = run_dir / "pytest_report.xml"
    coverage_data = run_dir / ".coverage"
    coverage_json = run_dir / "coverage.json"
    pytest_proc = run_pytest(
        args.python_exe,
        restored["workspace"],
        tests_file,
        restored["package_name"],
        junit_xml,
        coverage_data,
    )
    coverage_proc = run_coverage_json(
        args.python_exe,
        restored["workspace"],
        coverage_data,
        coverage_json,
    )

    junit_summary = parse_junit_xml(junit_xml)
    coverage_summary = parse_coverage(coverage_json, restored["target_file"])

    payload = {
        "model_name": args.model_name,
        "sample": {
            "id": sample["id"],
            "category": sample["category"],
            "package": sample["package"],
            "source_path": sample["source_path"],
            "module_import": restored["module_import"],
        },
        "workspace": str(restored["workspace"]),
        "generated_test_file": str(tests_file),
        "compile": {
            "passed": compile_ok,
            "output": compile_output,
        },
        "pytest": {
            "returncode": pytest_proc.returncode,
            "passed": pytest_proc.returncode == 0,
            "summary": junit_summary,
            "stdout": pytest_proc.stdout,
            "stderr": pytest_proc.stderr,
        },
        "coverage_command": {
            "returncode": coverage_proc.returncode,
            "stdout": coverage_proc.stdout,
            "stderr": coverage_proc.stderr,
        },
        "coverage": coverage_summary,
    }

    report_path = save_report(run_dir, payload)
    print(json.dumps(payload, ensure_ascii=False, indent=2))
    print(f"\nReport saved to: {report_path}")
    return 0 if pytest_proc.returncode == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
