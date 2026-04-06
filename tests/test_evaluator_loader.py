from __future__ import annotations

from pathlib import Path

from benchmark.evaluator.loader import load_generated_samples, load_generated_samples_with_mode


def test_loader_collects_generated_tests() -> None:
    samples = load_generated_samples(Path("results"))
    assert isinstance(samples, list)
    # 当前仓库通常已有历史生成结果；若为空也不应报错。
    for row in samples[:5]:
        assert row.generated_test_path.endswith((".test.py", ".test.java", ".test.go", ".test.cpp"))


def test_loader_keeps_latest_test_per_sample(tmp_path: Path) -> None:
    model_tests = tmp_path / "m1" / "tests"
    model_tests.mkdir(parents=True)

    old_file = model_tests / "m1_python_boundary_python_boundary_0_20260401010101.test.py"
    new_file = model_tests / "m1_python_boundary_python_boundary_0_20260402010101.test.py"
    other_file = model_tests / "m1_python_simple_function_python_simple_function_0_20260402010101.test.py"
    old_file.write_text("def test_old():\n    assert True\n", encoding="utf-8")
    new_file.write_text("def test_new():\n    assert True\n", encoding="utf-8")
    other_file.write_text("def test_other():\n    assert True\n", encoding="utf-8")

    samples = load_generated_samples(tmp_path)
    names = sorted(Path(s.generated_test_path).name for s in samples)
    assert names == sorted([new_file.name, other_file.name])


def test_loader_can_include_all_history(tmp_path: Path) -> None:
    model_tests = tmp_path / "m1" / "tests"
    model_tests.mkdir(parents=True)

    old_file = model_tests / "m1_python_boundary_python_boundary_0_20260401010101.test.py"
    new_file = model_tests / "m1_python_boundary_python_boundary_0_20260402010101.test.py"
    old_file.write_text("def test_old():\n    assert True\n", encoding="utf-8")
    new_file.write_text("def test_new():\n    assert True\n", encoding="utf-8")

    samples = load_generated_samples_with_mode(tmp_path, latest_only=False)
    names = sorted(Path(s.generated_test_path).name for s in samples)
    assert names == sorted([old_file.name, new_file.name])
