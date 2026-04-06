from __future__ import annotations

import re
from pathlib import Path

from .contracts import EvaluationSample


_TEST_FILE_PATTERN = re.compile(
    r"^(?P<model>.+?)_(?P<lang>python|java|go|cpp|javascript)_(?P<sample>.+?)_\d{14}\.test\.[a-zA-Z0-9]+$"
)


def load_generated_samples(results_root: Path) -> list[EvaluationSample]:
    return load_generated_samples_with_mode(results_root=results_root, latest_only=True)


def load_generated_samples_with_mode(
    results_root: Path,
    latest_only: bool,
) -> list[EvaluationSample]:
    """从 results/<model>/tests 中收集待评测样本。"""
    samples: list[EvaluationSample] = []
    if not results_root.exists():
        return samples

    for model_dir in sorted(results_root.iterdir()):
        if not model_dir.is_dir():
            continue
        tests_dir = model_dir / "tests"
        if not tests_dir.exists() or not tests_dir.is_dir():
            continue

        test_files: list[Path] = []
        if latest_only:
            latest_by_sample: dict[tuple[str, str], Path] = {}
            for test_file in sorted(tests_dir.iterdir()):
                if not test_file.is_file():
                    continue
                matched = _TEST_FILE_PATTERN.match(test_file.name)
                if not matched:
                    continue

                language = matched.group("lang")
                sample_id = matched.group("sample")
                key = (language, sample_id)
                prev = latest_by_sample.get(key)
                if prev is None or test_file.name > prev.name:
                    latest_by_sample[key] = test_file
            test_files = sorted(latest_by_sample.values())
        else:
            for test_file in sorted(tests_dir.iterdir()):
                if test_file.is_file() and _TEST_FILE_PATTERN.match(test_file.name):
                    test_files.append(test_file)

        for test_file in test_files:
            matched = _TEST_FILE_PATTERN.match(test_file.name)
            if not matched:
                continue

            language = matched.group("lang")
            sample_id = matched.group("sample")
            source_path = _find_source_path(sample_id=sample_id, language=language)
            samples.append(
                EvaluationSample(
                    model=matched.group("model"),
                    language=language,
                    sample_id=sample_id,
                    generated_test_path=str(test_file),
                    source_path=str(source_path) if source_path else None,
                    timestamp=test_file.stem,
                )
            )
    return samples


def _find_source_path(sample_id: str, language: str) -> Path | None:
    """根据 sample_id 在 dataset 下尝试匹配源码路径。"""
    dataset_dir = Path("dataset") / language
    if not dataset_dir.exists():
        return None

    for path in dataset_dir.rglob("*"):
        if not path.is_file():
            continue
        if path.stem == sample_id:
            return path
    return None
