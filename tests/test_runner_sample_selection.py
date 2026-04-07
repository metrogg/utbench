from __future__ import annotations

import tempfile
from pathlib import Path

from benchmark.runner.runner import Runner


def _make_runner(dataset_root: Path) -> Runner:
    runner = object.__new__(Runner)
    runner.dataset_root = dataset_root
    runner.results_root = dataset_root
    return runner


def test_collect_samples_supports_sample_glob_and_max_samples() -> None:
    with tempfile.TemporaryDirectory(prefix="runner_glob_") as td:
        root = Path(td)
        py_dir = root / "python"
        (py_dir / "boundary").mkdir(parents=True)
        (py_dir / "simple_function").mkdir(parents=True)
        (py_dir / "boundary" / "boundary_000.py").write_text("x=1\n", encoding="utf-8")
        (py_dir / "boundary" / "boundary_001.py").write_text("x=2\n", encoding="utf-8")
        (py_dir / "simple_function" / "simple_function_000.py").write_text(
            "x=3\n", encoding="utf-8"
        )

        runner = _make_runner(root)

        selected = Runner._collect_samples(
            runner,
            languages=["python"],
            max_samples=1,
            sample_glob="boundary/*.py",
        )
        assert len(selected) == 1
        assert selected[0].name == "boundary_000.py"

        selected_all = Runner._collect_samples(
            runner,
            languages=["python"],
            max_samples=None,
            sample_glob="boundary/*.py",
        )
        assert [path.name for path in selected_all] == ["boundary_000.py", "boundary_001.py"]


def test_infer_sample_language_prefers_dataset_top_level() -> None:
    with tempfile.TemporaryDirectory(prefix="runner_lang_") as td:
        root = Path(td)
        sample = root / "python" / "boundary" / "boundary_000.py"
        sample.parent.mkdir(parents=True)
        sample.write_text("x=1\n", encoding="utf-8")

        runner = _make_runner(root)
        runner.LANGUAGE_EXT = Runner.LANGUAGE_EXT

        language = Runner._infer_sample_language(runner, sample)
        assert language == "python"


def test_checkpoint_scope_includes_sample_glob() -> None:
    with tempfile.TemporaryDirectory(prefix="runner_ckpt_") as td:
        root = Path(td)
        runner = _make_runner(root)

        ckpt_a = Runner._checkpoint_path(
            runner,
            model_names=["deepseek"],
            languages=["python"],
            max_samples=5,
            sample_glob="boundary/*.py",
        )
        ckpt_b = Runner._checkpoint_path(
            runner,
            model_names=["deepseek"],
            languages=["python"],
            max_samples=5,
            sample_glob="simple_function/*.py",
        )

        assert ckpt_a != ckpt_b
