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
) -> tuple[bool | None, str | None, int | None]:
    """阶段2：执行测试。"""
    if language != "python":
        return None, f"Language test execution not implemented yet: {language}", None

    workdir, target_test = prepare_python_execution_workspace(
        generated_test_path=generated_test_path,
        source_path=source_path,
    )

    # 使用当前解释器，避免不同 shell 下 python 命令不可用。
    command = [sys.executable, "-m", "pytest", str(target_test.name), "-q"]
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
        if proc.returncode == 0:
            return True, None, elapsed_ms

        error = (proc.stdout + "\n" + proc.stderr).strip()
        return False, error[-4000:], elapsed_ms
    except subprocess.TimeoutExpired:
        elapsed_ms = int((time.perf_counter() - started) * 1000)
        return False, f"pytest timeout after {timeout_seconds}s", elapsed_ms
    finally:
        cleanup_execution_workspace(workdir)
