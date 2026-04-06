from __future__ import annotations

import ast
import json
import os
import signal
import subprocess
import sys
from pathlib import Path

from .compiler import prepare_python_execution_workspace, python_source_alias_candidates


def collect_coverage(
    language: str,
    generated_test_path: Path,
    source_path: Path | None,
) -> tuple[float | None, float | None, float | None, str | None]:
    """阶段3：覆盖率分析。

    Python 真实执行 coverage.py：
    - coverage run -m pytest
    - coverage json
    先输出 line coverage；branch/function 在 MVP 阶段先保留为 None。
    """
    if language == "python":
        return _collect_python_coverage(generated_test_path, source_path)

    if language in {"java", "go", "cpp", "javascript"}:
        return None, None, None, None

    return None, None, None, f"Unsupported language: {language}"


def _collect_python_coverage(
    generated_test_path: Path,
    source_path: Path | None,
) -> tuple[float | None, float | None, float | None, str | None]:
    if source_path is None or not source_path.exists():
        return None, None, None, "Missing source file for coverage"

    workdir, target_test = prepare_python_execution_workspace(
        generated_test_path=generated_test_path,
        source_path=source_path,
    )

    source_name = source_path.name
    source_aliases = python_source_alias_candidates(
        source_path=source_path,
        generated_test_path=generated_test_path,
    )
    source_aliases.add(source_name)
    json_report = workdir / "coverage.json"

    run_cmd = [
        sys.executable,
        "-m",
        "coverage",
        "run",
        "--branch",
        "-m",
        "pytest",
        str(target_test.name),
        "-q",
    ]
    proc = subprocess.run(
        run_cmd,
        cwd=str(workdir),
        capture_output=True,
        text=True,
        timeout=40,
        check=False,
    )
    if proc.returncode != 0:
        error = (proc.stdout + "\n" + proc.stderr).strip()
        return None, None, None, f"coverage run failed: {error[-3000:]}"

    json_cmd = [
        sys.executable,
        "-m",
        "coverage",
        "json",
        "-o",
        str(json_report),
    ]
    proc_json = subprocess.run(
        json_cmd,
        cwd=str(workdir),
        capture_output=True,
        text=True,
        timeout=20,
        check=False,
    )
    if proc_json.returncode != 0:
        error = (proc_json.stdout + "\n" + proc_json.stderr).strip()
        return None, None, None, f"coverage json failed: {error[-3000:]}"

    if not json_report.exists():
        return None, None, None, "coverage json report not found"

    payload = json.loads(json_report.read_text(encoding="utf-8"))
    files = payload.get("files", {})

    line_coverage = None
    branch_coverage = None
    # 优先用被测源码文件对应的覆盖率。
    for file_path, details in files.items():
        file_name = Path(file_path).name
        file_stem = Path(file_path).stem
        if file_name in source_aliases or file_stem in source_aliases:
            summary = details.get("summary", {})
            line_coverage = _extract_line_coverage(summary)
            branch_coverage = _extract_branch_coverage(summary)
            break

    if line_coverage is None:
        return None, None, None, (
            f"source file not found in coverage report: {source_name}; aliases={sorted(source_aliases)}"
        )

    # 源码命中后，若分支覆盖率缺失再回退到 totals（兼容不同 coverage 版本输出差异）。
    if branch_coverage is None:
        totals = payload.get("totals", {})
        branch_coverage = _extract_branch_coverage(totals)

    return line_coverage, branch_coverage, None, None


def _extract_line_coverage(summary: dict) -> float | None:
    """提取行覆盖率并统一到 0~1。"""
    if summary.get("covered_lines") is not None and summary.get("num_statements"):
        covered = float(summary.get("covered_lines", 0))
        total = float(summary.get("num_statements", 0))
        if total > 0:
            return round(covered / total, 6)

    percent = summary.get("percent_covered")
    if percent is not None:
        return round(float(percent) / 100.0, 6)
    return None


def _extract_branch_coverage(summary: dict) -> float | None:
    """提取分支覆盖率并统一到 0~1。"""
    covered_branches = summary.get("covered_branches")
    num_branches = summary.get("num_branches")
    if covered_branches is not None and num_branches is not None:
        total = float(num_branches)
        if total > 0:
            return round(float(covered_branches) / total, 6)
        # 无分支时返回 None，避免误导性 100%。
        return None

    percent_branches = summary.get("percent_covered_branches")
    if percent_branches is not None:
        return round(float(percent_branches) / 100.0, 6)
    return None


def collect_mutation_score(
    language: str,
    generated_test_path: Path,
    source_path: Path | None,
) -> tuple[float | None, str | None, dict[str, int] | None]:
    """阶段4：变异测试。当前 Python 使用 mutmut。"""
    if language == "python":
        return _collect_python_mutation_score(generated_test_path, source_path)

    if language in {"java", "go", "cpp", "javascript"}:
        return None, None, None

    return None, f"Unsupported language: {language}", None


def _collect_python_mutation_score(
    generated_test_path: Path,
    source_path: Path | None,
) -> tuple[float | None, str | None, dict[str, int] | None]:
    if source_path is None or not source_path.exists():
        return None, "Missing source file for mutation", None

    workdir, target_test = prepare_python_execution_workspace(
        generated_test_path=generated_test_path,
        source_path=source_path,
    )

    mutation_targets = _infer_mutation_targets(
        workdir=workdir,
        target_test=target_test,
        fallback_source=source_path.name,
    )

    pyproject_file = workdir / "pyproject.toml"
    pyproject_file.write_text(
        _build_mutmut_pyproject(source_names=mutation_targets, test_name=target_test.name),
        encoding="utf-8",
    )
    mutmut_env = _build_mutmut_env(workdir)

    run_cmd = [sys.executable, "-m", "mutmut", "run", "--max-children", "1"]
    run_proc = subprocess.run(
        run_cmd,
        cwd=str(workdir),
        env=mutmut_env,
        capture_output=True,
        text=True,
        timeout=120,
        check=False,
    )

    export_cmd = [sys.executable, "-m", "mutmut", "export-cicd-stats"]
    export_proc = subprocess.run(
        export_cmd,
        cwd=str(workdir),
        env=mutmut_env,
        capture_output=True,
        text=True,
        timeout=30,
        check=False,
    )

    stats_file = workdir / "mutants" / "mutmut-cicd-stats.json"
    if not stats_file.exists():
        return None, _format_mutmut_error(
            prefix="mutmut stats file not found",
            run_proc=run_proc,
            export_proc=export_proc,
        ), None

    try:
        payload = json.loads(stats_file.read_text(encoding="utf-8"))
    except Exception as exc:  # pylint: disable=broad-except
        return None, f"mutmut stats parse failed: {exc}", None

    stats = _extract_mutation_stats(payload)
    if stats is None:
        return None, f"invalid mutmut stats payload: {payload}", None

    meta_stats = _extract_meta_mutation_stats(workdir, mutation_targets)
    if meta_stats is not None:
        stats = meta_stats

    killed = _to_int(payload.get("killed"))
    total = _to_int(payload.get("total"))
    if killed is None or total is None:
        return None, f"invalid mutmut stats payload: {payload}", stats
    if total <= 0:
        return None, "mutmut produced zero mutants", stats

    processed = _processed_mutants(stats)
    if processed <= 0:
        return None, _format_mutmut_error(
            prefix="mutmut did not execute any mutants",
            run_proc=run_proc,
            export_proc=export_proc,
        ), stats
    if (run_proc.returncode != 0 or stats.get("not_checked", 0) > 0) and processed < total:
        return None, _format_mutmut_error(
            prefix=f"mutmut run incomplete ({processed}/{total})",
            run_proc=run_proc,
            export_proc=export_proc,
        ), stats

    score = round(killed / total, 6)
    return score, None, stats


def _build_mutmut_pyproject(source_names: list[str], test_name: str) -> str:
    source_literal = json.dumps(source_names)
    test_literal = json.dumps(test_name)
    return (
        "[tool.mutmut]\n"
        f"paths_to_mutate = {source_literal}\n"
        f"pytest_add_cli_args_test_selection = [{test_literal}]\n"
        "pytest_add_cli_args = [\"-q\"]\n"
    )


def _build_mutmut_env(workdir: Path) -> dict[str, str]:
    env = os.environ.copy()
    shim_dir = workdir / ".mutmut_shim"
    shim_dir.mkdir(parents=True, exist_ok=True)
    sitecustomize = shim_dir / "sitecustomize.py"
    sitecustomize.write_text(
        (
            "import multiprocessing as _mp\n"
            "import multiprocessing.context as _mpc\n"
            "_orig_set_start_method = _mpc._default_context.set_start_method\n"
            "def _safe_set_start_method(method, force=False):\n"
            "    try:\n"
            "        return _orig_set_start_method(method, force=force)\n"
            "    except RuntimeError as exc:\n"
            "        if 'context has already been set' in str(exc):\n"
            "            return None\n"
            "        raise\n"
            "_mp.set_start_method = _safe_set_start_method\n"
            "_mpc.set_start_method = _safe_set_start_method\n"
        ),
        encoding="utf-8",
    )

    current_path = env.get("PYTHONPATH", "")
    env["PYTHONPATH"] = (
        f"{shim_dir}{os.pathsep}{current_path}" if current_path else str(shim_dir)
    )
    return env


def _infer_mutation_targets(
    workdir: Path,
    target_test: Path,
    fallback_source: str,
) -> list[str]:
    imported_local_modules = _extract_local_imported_modules(workdir=workdir, target_test=target_test)
    if imported_local_modules:
        return sorted(f"{module}.py" for module in imported_local_modules)
    return [fallback_source]


def _extract_local_imported_modules(workdir: Path, target_test: Path) -> set[str]:
    modules: set[str] = set()
    try:
        tree = ast.parse(target_test.read_text(encoding="utf-8"))
    except Exception:  # pylint: disable=broad-except
        return modules

    for node in ast.walk(tree):
        if isinstance(node, ast.Import):
            for alias in node.names:
                top = alias.name.split(".", maxsplit=1)[0]
                if (workdir / f"{top}.py").exists():
                    modules.add(top)
        elif isinstance(node, ast.ImportFrom):
            if node.level and node.level > 0:
                continue
            if node.module:
                top = node.module.split(".", maxsplit=1)[0]
                if (workdir / f"{top}.py").exists():
                    modules.add(top)
    return modules


def _to_int(value: object) -> int | None:
    if isinstance(value, bool):
        return None
    if isinstance(value, int):
        return value
    if isinstance(value, float):
        return int(value)
    if isinstance(value, str):
        text = value.strip()
        if text and text.lstrip("-").isdigit():
            return int(text)
    return None


def _extract_mutation_stats(payload: dict[str, object]) -> dict[str, int] | None:
    required_total = _to_int(payload.get("total"))
    required_killed = _to_int(payload.get("killed"))
    if required_total is None or required_killed is None:
        return None

    def _value_or_zero(key: str) -> int:
        value = _to_int(payload.get(key))
        return value if value is not None else 0

    return {
        "total": required_total,
        "killed": required_killed,
        "not_checked": _value_or_zero("not_checked"),
        "survived": _value_or_zero("survived"),
        "timeout": _value_or_zero("timeout"),
        "no_tests": _value_or_zero("no_tests"),
        "skipped": _value_or_zero("skipped"),
        "suspicious": _value_or_zero("suspicious"),
    }


def _extract_meta_mutation_stats(
    workdir: Path,
    mutation_targets: list[str],
) -> dict[str, int] | None:
    mutants_dir = workdir / "mutants"
    if not mutants_dir.exists():
        return None

    stats = {
        "total": 0,
        "killed": 0,
        "survived": 0,
        "not_checked": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
        "timeout": 0,
    }
    has_any = False

    for target in mutation_targets:
        meta_file = mutants_dir / f"{target}.meta"
        if not meta_file.exists():
            continue
        has_any = True
        try:
            payload = json.loads(meta_file.read_text(encoding="utf-8"))
        except Exception:  # pylint: disable=broad-except
            continue

        exit_code_by_key = payload.get("exit_code_by_key", {})
        if not isinstance(exit_code_by_key, dict):
            continue

        for code in exit_code_by_key.values():
            stats["total"] += 1
            status = _mutmut_status_from_exit_code(code)
            if status in stats:
                stats[status] += 1

    if not has_any or stats["total"] <= 0:
        return None
    return stats


def _mutmut_status_from_exit_code(code: object) -> str:
    if code is None:
        return "not_checked"
    parsed = _to_int(code)
    if parsed is None:
        return "suspicious"

    mapping = {
        1: "killed",
        3: "killed",
        -24: "killed",
        0: "survived",
        5: "no_tests",
        2: "check_was_interrupted_by_user",
        33: "no_tests",
        -11: "segfault",
        -signal.SIGXCPU: "timeout",
    }
    return mapping.get(parsed, "suspicious")


def _processed_mutants(stats: dict[str, int]) -> int:
    fields = (
        "killed",
        "survived",
        "no_tests",
        "skipped",
        "suspicious",
        "timeout",
    )
    return sum(stats.get(field, 0) for field in fields)


def _format_mutmut_error(
    prefix: str,
    run_proc: subprocess.CompletedProcess[str],
    export_proc: subprocess.CompletedProcess[str],
) -> str:
    run_stdout = (run_proc.stdout or "").strip()
    run_stderr = (run_proc.stderr or "").strip()
    export_stdout = (export_proc.stdout or "").strip()
    export_stderr = (export_proc.stderr or "").strip()
    run_stdout_snippet = run_stdout[-1200:] if run_stdout else ""
    run_snippet = run_stderr[-1200:] if run_stderr else ""
    export_stdout_snippet = export_stdout[-1200:] if export_stdout else ""
    export_snippet = export_stderr[-1200:] if export_stderr else ""
    return (
        f"{prefix}; "
        f"mutmut run rc={run_proc.returncode}; "
        f"mutmut export rc={export_proc.returncode}; "
        f"run_stdout={run_stdout_snippet!r}; "
        f"run_stderr={run_snippet!r}; "
        f"export_stdout={export_stdout_snippet!r}; "
        f"export_stderr={export_snippet!r}"
    )
