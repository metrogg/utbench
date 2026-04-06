from __future__ import annotations

import json
import subprocess
from pathlib import Path

from benchmark.evaluator import coverage as coverage_mod


def test_collect_mutation_score_python_reads_mutmut_stats(
    tmp_path: Path,
    monkeypatch,
) -> None:
    src = tmp_path / "calculator.py"
    src.write_text("def add(a, b):\n    return a + b\n", encoding="utf-8")

    generated = tmp_path / "generated.test.py"
    generated.write_text(
        "from calculator import add\n\n"
        "def test_add():\n"
        "    assert add(1, 2) == 3\n",
        encoding="utf-8",
    )

    workdir = tmp_path / "workdir"
    workdir.mkdir()
    target_test = workdir / "test_generated.py"
    target_test.write_text(generated.read_text(encoding="utf-8"), encoding="utf-8")

    monkeypatch.setattr(
        coverage_mod,
        "prepare_python_execution_workspace",
        lambda generated_test_path, source_path: (workdir, target_test),
    )

    def fake_run(command, **kwargs):  # type: ignore[no-untyped-def]
        _ = kwargs
        if command[:4] == ["python3", "-m", "mutmut", "run"]:
            return subprocess.CompletedProcess(command, 0, "", "")
        if command[:4] == ["python3", "-m", "mutmut", "export-cicd-stats"]:
            stats_dir = workdir / "mutants"
            stats_dir.mkdir(parents=True, exist_ok=True)
            (stats_dir / "mutmut-cicd-stats.json").write_text(
                json.dumps({"killed": 3, "total": 5}),
                encoding="utf-8",
            )
            return subprocess.CompletedProcess(command, 0, "", "")
        return subprocess.CompletedProcess(command, 0, "", "")

    monkeypatch.setattr(coverage_mod.subprocess, "run", fake_run)
    monkeypatch.setattr(coverage_mod.sys, "executable", "python3")

    score, err, stats = coverage_mod.collect_mutation_score(
        language="python",
        generated_test_path=generated,
        source_path=src,
    )

    assert err is None
    assert score == 0.6
    assert stats == {
        "total": 5,
        "killed": 3,
        "not_checked": 0,
        "survived": 0,
        "timeout": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
    }
    pyproject_text = (workdir / "pyproject.toml").read_text(encoding="utf-8")
    assert "paths_to_mutate = [\"calculator.py\"]" in pyproject_text
    assert "pytest_add_cli_args_test_selection = [\"test_generated.py\"]" in pyproject_text


def test_collect_mutation_score_python_targets_imported_alias_module(
    tmp_path: Path,
    monkeypatch,
) -> None:
    src = tmp_path / "boundary_python_boundary_0.py"
    src.write_text("def below_zero(ops):\n    return False\n", encoding="utf-8")
    generated = tmp_path / "generated.test.py"
    generated.write_text(
        "from solution import below_zero\n\n"
        "def test_stub():\n"
        "    assert below_zero([1, 2, 3]) is False\n",
        encoding="utf-8",
    )

    workdir = tmp_path / "work_alias"
    workdir.mkdir()
    target_test = workdir / "test_generated.py"
    target_test.write_text(generated.read_text(encoding="utf-8"), encoding="utf-8")
    (workdir / "boundary_python_boundary_0.py").write_text(
        src.read_text(encoding="utf-8"),
        encoding="utf-8",
    )
    (workdir / "solution.py").write_text(src.read_text(encoding="utf-8"), encoding="utf-8")

    monkeypatch.setattr(
        coverage_mod,
        "prepare_python_execution_workspace",
        lambda generated_test_path, source_path: (workdir, target_test),
    )

    def fake_run(command, **kwargs):  # type: ignore[no-untyped-def]
        _ = kwargs
        if command[:4] == ["python3", "-m", "mutmut", "export-cicd-stats"]:
            stats_dir = workdir / "mutants"
            stats_dir.mkdir(parents=True, exist_ok=True)
            (stats_dir / "mutmut-cicd-stats.json").write_text(
                json.dumps({"killed": 1, "total": 2}),
                encoding="utf-8",
            )
        return subprocess.CompletedProcess(command, 0, "", "")

    monkeypatch.setattr(coverage_mod.subprocess, "run", fake_run)
    monkeypatch.setattr(coverage_mod.sys, "executable", "python3")

    score, err, stats = coverage_mod.collect_mutation_score(
        language="python",
        generated_test_path=generated,
        source_path=src,
    )

    assert err is None
    assert score == 0.5
    assert stats == {
        "total": 2,
        "killed": 1,
        "not_checked": 0,
        "survived": 0,
        "timeout": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
    }
    pyproject_text = (workdir / "pyproject.toml").read_text(encoding="utf-8")
    assert "paths_to_mutate = [\"solution.py\"]" in pyproject_text


def test_build_mutmut_env_creates_sitecustomize(tmp_path: Path) -> None:
    from benchmark.evaluator.coverage import _build_mutmut_env

    env = _build_mutmut_env(tmp_path)
    shim = tmp_path / ".mutmut_shim" / "sitecustomize.py"
    assert shim.exists()
    assert "PYTHONPATH" in env
    assert str(tmp_path / ".mutmut_shim") in env["PYTHONPATH"]


def test_collect_mutation_score_python_handles_zero_mutants(
    tmp_path: Path,
    monkeypatch,
) -> None:
    src = tmp_path / "calc.py"
    src.write_text("def inc(v):\n    return v + 1\n", encoding="utf-8")
    generated = tmp_path / "generated.py"
    generated.write_text("def test_placeholder():\n    assert True\n", encoding="utf-8")

    workdir = tmp_path / "work"
    workdir.mkdir()
    target_test = workdir / "test_generated.py"
    target_test.write_text(generated.read_text(encoding="utf-8"), encoding="utf-8")

    monkeypatch.setattr(
        coverage_mod,
        "prepare_python_execution_workspace",
        lambda generated_test_path, source_path: (workdir, target_test),
    )

    def fake_run(command, **kwargs):  # type: ignore[no-untyped-def]
        _ = kwargs
        if command[:4] == ["python3", "-m", "mutmut", "export-cicd-stats"]:
            stats_dir = workdir / "mutants"
            stats_dir.mkdir(parents=True, exist_ok=True)
            (stats_dir / "mutmut-cicd-stats.json").write_text(
                json.dumps({"killed": 0, "total": 0}),
                encoding="utf-8",
            )
        return subprocess.CompletedProcess(command, 0, "", "")

    monkeypatch.setattr(coverage_mod.subprocess, "run", fake_run)
    monkeypatch.setattr(coverage_mod.sys, "executable", "python3")

    score, err, stats = coverage_mod.collect_mutation_score(
        language="python",
        generated_test_path=generated,
        source_path=src,
    )

    assert score is None
    assert err == "mutmut produced zero mutants"
    assert stats == {
        "total": 0,
        "killed": 0,
        "not_checked": 0,
        "survived": 0,
        "timeout": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
    }


def test_collect_mutation_score_python_fails_when_no_mutants_executed(
    tmp_path: Path,
    monkeypatch,
) -> None:
    src = tmp_path / "calc.py"
    src.write_text("def inc(v):\n    return v + 1\n", encoding="utf-8")
    generated = tmp_path / "generated.py"
    generated.write_text("from calc import inc\n\ndef test_inc():\n    assert inc(1) == 2\n", encoding="utf-8")

    workdir = tmp_path / "work_no_exec"
    workdir.mkdir()
    target_test = workdir / "test_generated.py"
    target_test.write_text(generated.read_text(encoding="utf-8"), encoding="utf-8")
    (workdir / "calc.py").write_text(src.read_text(encoding="utf-8"), encoding="utf-8")

    monkeypatch.setattr(
        coverage_mod,
        "prepare_python_execution_workspace",
        lambda generated_test_path, source_path: (workdir, target_test),
    )

    def fake_run(command, **kwargs):  # type: ignore[no-untyped-def]
        _ = kwargs
        if command[:4] == ["python3", "-m", "mutmut", "run"]:
            return subprocess.CompletedProcess(command, 1, "", "runtime issue")
        if command[:4] == ["python3", "-m", "mutmut", "export-cicd-stats"]:
            stats_dir = workdir / "mutants"
            stats_dir.mkdir(parents=True, exist_ok=True)
            (stats_dir / "mutmut-cicd-stats.json").write_text(
                json.dumps({"killed": 0, "survived": 0, "total": 8}),
                encoding="utf-8",
            )
        return subprocess.CompletedProcess(command, 0, "", "")

    monkeypatch.setattr(coverage_mod.subprocess, "run", fake_run)
    monkeypatch.setattr(coverage_mod.sys, "executable", "python3")

    score, err, stats = coverage_mod.collect_mutation_score(
        language="python",
        generated_test_path=generated,
        source_path=src,
    )

    assert score is None
    assert err is not None
    assert err.startswith("mutmut did not execute any mutants")
    assert stats == {
        "total": 8,
        "killed": 0,
        "not_checked": 0,
        "survived": 0,
        "timeout": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
    }


def test_collect_mutation_score_python_uses_meta_not_checked_counts(
    tmp_path: Path,
    monkeypatch,
) -> None:
    src = tmp_path / "boundary.py"
    src.write_text("def below_zero(ops):\n    return False\n", encoding="utf-8")
    generated = tmp_path / "generated.py"
    generated.write_text("def test_placeholder():\n    assert True\n", encoding="utf-8")

    workdir = tmp_path / "work_meta"
    workdir.mkdir()
    target_test = workdir / "test_generated.py"
    target_test.write_text(generated.read_text(encoding="utf-8"), encoding="utf-8")
    (workdir / "boundary.py").write_text(src.read_text(encoding="utf-8"), encoding="utf-8")

    monkeypatch.setattr(
        coverage_mod,
        "prepare_python_execution_workspace",
        lambda generated_test_path, source_path: (workdir, target_test),
    )
    monkeypatch.setattr(
        coverage_mod,
        "_infer_mutation_targets",
        lambda workdir, target_test, fallback_source: ["boundary.py"],
    )

    def fake_run(command, **kwargs):  # type: ignore[no-untyped-def]
        _ = kwargs
        if command[:4] == ["python3", "-m", "mutmut", "run"]:
            mutants_dir = workdir / "mutants"
            mutants_dir.mkdir(parents=True, exist_ok=True)
            (mutants_dir / "boundary.py.meta").write_text(
                json.dumps(
                    {
                        "exit_code_by_key": {
                            "boundary.x_1": None,
                            "boundary.x_2": None,
                            "boundary.x_3": None,
                            "boundary.x_4": None,
                        }
                    }
                ),
                encoding="utf-8",
            )
            return subprocess.CompletedProcess(command, 1, "", "")
        if command[:4] == ["python3", "-m", "mutmut", "export-cicd-stats"]:
            stats_dir = workdir / "mutants"
            stats_dir.mkdir(parents=True, exist_ok=True)
            (stats_dir / "mutmut-cicd-stats.json").write_text(
                json.dumps({"killed": 0, "survived": 0, "total": 4}),
                encoding="utf-8",
            )
        return subprocess.CompletedProcess(command, 0, "", "")

    monkeypatch.setattr(coverage_mod.subprocess, "run", fake_run)
    monkeypatch.setattr(coverage_mod.sys, "executable", "python3")

    score, err, stats = coverage_mod.collect_mutation_score(
        language="python",
        generated_test_path=generated,
        source_path=src,
    )

    assert score is None
    assert err is not None
    assert err.startswith("mutmut did not execute any mutants")
    assert stats == {
        "total": 4,
        "killed": 0,
        "not_checked": 4,
        "survived": 0,
        "timeout": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
    }


def test_collect_mutation_score_python_reports_incomplete_when_not_checked(
    tmp_path: Path,
    monkeypatch,
) -> None:
    src = tmp_path / "sample.py"
    src.write_text("def f(x):\n    return x + 1\n", encoding="utf-8")
    generated = tmp_path / "generated.py"
    generated.write_text("def test_placeholder():\n    assert True\n", encoding="utf-8")

    workdir = tmp_path / "work_incomplete"
    workdir.mkdir()
    target_test = workdir / "test_generated.py"
    target_test.write_text(generated.read_text(encoding="utf-8"), encoding="utf-8")
    (workdir / "sample.py").write_text(src.read_text(encoding="utf-8"), encoding="utf-8")

    monkeypatch.setattr(
        coverage_mod,
        "prepare_python_execution_workspace",
        lambda generated_test_path, source_path: (workdir, target_test),
    )
    monkeypatch.setattr(
        coverage_mod,
        "_infer_mutation_targets",
        lambda workdir, target_test, fallback_source: ["sample.py"],
    )

    def fake_run(command, **kwargs):  # type: ignore[no-untyped-def]
        _ = kwargs
        if command[:4] == ["python3", "-m", "mutmut", "run"]:
            mutants_dir = workdir / "mutants"
            mutants_dir.mkdir(parents=True, exist_ok=True)
            (mutants_dir / "sample.py.meta").write_text(
                json.dumps(
                    {
                        "exit_code_by_key": {
                            "sample.x_1": 1,
                            "sample.x_2": None,
                            "sample.x_3": None,
                        }
                    }
                ),
                encoding="utf-8",
            )
            return subprocess.CompletedProcess(command, 1, "", "")
        if command[:4] == ["python3", "-m", "mutmut", "export-cicd-stats"]:
            stats_dir = workdir / "mutants"
            stats_dir.mkdir(parents=True, exist_ok=True)
            (stats_dir / "mutmut-cicd-stats.json").write_text(
                json.dumps({"killed": 1, "survived": 0, "total": 3}),
                encoding="utf-8",
            )
        return subprocess.CompletedProcess(command, 0, "", "")

    monkeypatch.setattr(coverage_mod.subprocess, "run", fake_run)
    monkeypatch.setattr(coverage_mod.sys, "executable", "python3")

    score, err, stats = coverage_mod.collect_mutation_score(
        language="python",
        generated_test_path=generated,
        source_path=src,
    )

    assert score is None
    assert err is not None
    assert err.startswith("mutmut run incomplete")
    assert stats == {
        "total": 3,
        "killed": 1,
        "not_checked": 2,
        "survived": 0,
        "timeout": 0,
        "no_tests": 0,
        "skipped": 0,
        "suspicious": 0,
    }
