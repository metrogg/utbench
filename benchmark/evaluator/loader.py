from __future__ import annotations

import json
from pathlib import Path

from .contracts import EvaluationSample


_LANGUAGE_ALIASES = {
    "python": "python",
    "py": "python",
    "java": "java",
    "go": "go",
    "cpp": "cpp",
    "c++": "cpp",
    "javascript": "javascript",
    "js": "javascript",
    # Legacy runner bug: scenario folder was used as language.
    "boundary": "python",
    "simple_function": "python",
    "complex_dependency": "python",
    "interface_mock": "python",
}


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
                parsed = _parse_test_file_name(model=model_dir.name, file_name=test_file.name)
                if parsed is None:
                    continue

                language, sample_id = parsed
                key = (language, sample_id)
                prev = latest_by_sample.get(key)
                if prev is None or test_file.name > prev.name:
                    latest_by_sample[key] = test_file
            test_files = sorted(latest_by_sample.values())
        else:
            for test_file in sorted(tests_dir.iterdir()):
                parsed = _parse_test_file_name(model=model_dir.name, file_name=test_file.name)
                if test_file.is_file() and parsed is not None:
                    test_files.append(test_file)

        for test_file in test_files:
            parsed = _parse_test_file_name(model=model_dir.name, file_name=test_file.name)
            if parsed is None:
                continue

            language, sample_id = parsed
            source_path = _resolve_source_path(
                test_file=test_file,
                sample_id=sample_id,
                language=language,
            )
            samples.append(
                EvaluationSample(
                    model=model_dir.name,
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

    ext_map = {
        "python": ".py",
        "java": ".java",
        "go": ".go",
        "cpp": ".cpp",
        "javascript": ".js",
    }
    expected_ext = ext_map.get(language)
    if expected_ext:
        direct = dataset_dir / f"{sample_id}{expected_ext}"
        if direct.exists() and direct.is_file():
            return direct

    for path in dataset_dir.rglob("*"):
        if not path.is_file():
            continue
        if path.stem == sample_id:
            return path
    return None


def _parse_test_file_name(model: str, file_name: str) -> tuple[str, str] | None:
    model_prefix = f"{model}_"
    if not file_name.startswith(model_prefix):
        return None

    if ".test." not in file_name:
        return None

    payload = file_name[len(model_prefix) :]
    base = payload.split(".test.", maxsplit=1)[0]
    try:
        lang_and_sample, timestamp = base.rsplit("_", maxsplit=1)
    except ValueError:
        return None

    if len(timestamp) != 14 or not timestamp.isdigit():
        return None

    # Prefer longest alias first to support legacy tokens like "complex_dependency".
    aliases = sorted(_LANGUAGE_ALIASES.keys(), key=len, reverse=True)
    for alias in aliases:
        prefix = f"{alias}_"
        if not lang_and_sample.startswith(prefix):
            continue
        sample_id = lang_and_sample[len(prefix) :]
        if not sample_id:
            continue
        language = _normalize_language(alias, sample_id)
        if language:
            return language, sample_id

    # Fallback for unknown token shapes.
    try:
        language_token, sample_id = lang_and_sample.split("_", maxsplit=1)
    except ValueError:
        return None
    language = _normalize_language(language_token, sample_id)
    if language:
        return language, sample_id
    return None


def _normalize_language(raw_language: str, sample_id: str) -> str | None:
    lowered = raw_language.strip().lower()
    normalized = _LANGUAGE_ALIASES.get(lowered)
    if normalized:
        return normalized

    # Fallback for malformed legacy artifacts where language token is missing.
    sample_lower = sample_id.lower()
    if "_python_" in f"_{sample_lower}_":
        return "python"
    if "_java_" in f"_{sample_lower}_":
        return "java"
    if "_go_" in f"_{sample_lower}_":
        return "go"
    if "_javascript_" in f"_{sample_lower}_":
        return "javascript"
    return None


def _resolve_source_path(
    test_file: Path,
    sample_id: str,
    language: str,
) -> Path | None:
    metadata_path = test_file.parent.parent / "reports" / f"{test_file.stem.split('.test', 1)[0]}.metadata.json"
    if metadata_path.exists():
        try:
            payload = json.loads(metadata_path.read_text(encoding="utf-8"))
            sample_path = payload.get("sample_path")
            if isinstance(sample_path, str) and sample_path.strip():
                candidate = Path(sample_path)
                if candidate.exists() and candidate.is_file():
                    return candidate
        except Exception:
            pass

    return _find_source_path(sample_id=sample_id, language=language)
