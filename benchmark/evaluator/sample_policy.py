from __future__ import annotations

import fnmatch
from pathlib import Path


DEFAULT_ALLOWLIST_PATH = Path("benchmark/config/python_self_contained_allowlist.txt")


class PythonSelfContainedPolicy:
    """Classify Python samples by self-contained allowlist."""

    def __init__(self, allowlist_path: str | Path = DEFAULT_ALLOWLIST_PATH) -> None:
        self.allowlist_path = Path(allowlist_path)
        self._patterns = self._load_patterns()

    def classify(self, source_path: Path | None) -> tuple[str, str | None]:
        if source_path is None:
            return "unknown", "missing source file"

        rel_path = _dataset_python_relative_path(source_path)
        if rel_path is None:
            return "unknown", "outside dataset/python"

        if not self._patterns:
            return "unknown", "allowlist is empty"

        for pattern in self._patterns:
            if fnmatch.fnmatchcase(rel_path, pattern):
                return "self_contained", None

        return "non_self_contained", f"not allowlisted: {rel_path}"

    def _load_patterns(self) -> list[str]:
        if not self.allowlist_path.exists():
            return []

        patterns: list[str] = []
        for raw_line in self.allowlist_path.read_text(encoding="utf-8").splitlines():
            line = raw_line.strip()
            if not line or line.startswith("#"):
                continue
            patterns.append(line)
        return patterns


def _dataset_python_relative_path(source_path: Path) -> str | None:
    text = source_path.as_posix()
    marker = "/dataset/python/"
    if marker in text:
        rel = text.split(marker, maxsplit=1)[1]
        return rel or None
    if text.startswith("dataset/python/"):
        rel = text[len("dataset/python/") :]
        return rel or None
    return None
