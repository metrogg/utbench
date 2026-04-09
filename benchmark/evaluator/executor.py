from __future__ import annotations

import subprocess
import sys
import time
from pathlib import Path

from .compiler import cleanup_execution_workspace, prepare_python_execution_workspace


def execute_tests(
    language: str,
    generated_test_path: Path,
    source_path: Path | None,
    timeout_seconds: int = 30,
) -> tuple[bool | None, str | None, int | None, int | None, int | None]:
    """阶段2：执行测试。

    Returns:
        (test_pass, test_error, runtime_ms, passed_count, total_count)
        - test_pass: True if all passed, False if any failed, None if error.
        - passed_count: number of individual test functions passed.
        - total_count: total number of individual test functions.
    """
    if language != "python":
        return None, f"Language test execution not implemented yet: {language}", None, None, None

    workdir, target_test = prepare_python_execution_workspace(
        generated_test_path=generated_test_path,
        source_path=source_path,
    )

    command = [sys.executable, "-m", "pytest", str(target_test.name), "-q", "--maxfail=9999"]
    started = time.perf_counter()
    try:
        proc = subprocess.run(
            command,
            cwd=str(workdir),
            capture_output=True,
            text=True,
            timeout=timeout_seconds,
            check=False,
        )
        elapsed_ms = int((time.perf_counter() - started) * 1000)
        passed_count, total_count = _parse_pytest_counts(proc.stdout + "\n" + proc.stderr)
        if proc.returncode == 0:
            return True, None, elapsed_ms, passed_count, total_count

        error = (proc.stdout + "\n" + proc.stderr).strip()
        return False, error[-4000:], elapsed_ms, passed_count, total_count
    except subprocess.TimeoutExpired:
        elapsed_ms = int((time.perf_counter() - started) * 1000)
        return False, f"pytest timeout after {timeout_seconds}s", elapsed_ms, None, None
    finally:
        cleanup_execution_workspace(workdir)


def _parse_pytest_counts(output: str) -> tuple[int | None, int | None]:
    """从 pytest 输出解析 passed/total 数量。"""
    import re
    passed = None
    total = None
    m = re.search(r"(\d+) passed", output)
    if m:
        passed = int(m.group(1))
    m = re.search(r"(\d+) failed", output)
    if m:
        failed = int(m.group(1))
        if passed is not None:
            total = passed + failed
        elif failed > 0:
            total = None
    if passed is not None and total is None:
        total = passed
    return passed, total
