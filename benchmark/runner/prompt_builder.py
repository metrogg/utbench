from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path
from typing import Any


@dataclass
class PromptInput:
    # 样本 ID（来自文件名）
    sample_id: str
    # 编程语言标识
    language: str
    # 待测源码内容
    source_code: str
    # 样本文件路径
    sample_path: Path


class PromptBuilder:
    """Build prompt text from dataset sample and benchmark requirements."""

    LANGUAGE_FRAMEWORK = {
        "java": "JUnit 4",
        "python": "pytest",
        "go": "Go testing package",
        "cpp": "GoogleTest",
        "javascript": "Jest",
    }

    def __init__(self, benchmark_config: dict[str, Any] | None = None) -> None:
        benchmark_config = benchmark_config or {}
        self._thresholds = benchmark_config.get("thresholds", {}) or {}

    def read_sample(
        self,
        sample_path: str | Path,
        default_language: str | None = None,
    ) -> PromptInput:
        sample = Path(sample_path)
        source_code = sample.read_text(encoding="utf-8")
        language = self._infer_language(sample, default_language)
        return PromptInput(
            sample_id=sample.stem,
            language=language,
            source_code=source_code,
            sample_path=sample,
        )

    def build_prompt(self, data: PromptInput) -> str:
        framework = self.LANGUAGE_FRAMEWORK.get(data.language, "the standard test framework")
        dependencies = self._extract_dependencies(data.source_code, data.language)
        scenario, complexity = self._parse_sample_meta(data.sample_id, data.sample_path)
        module_name = self._module_import_name(data.sample_id)
        coverage_targets = self._coverage_targets_text()
        mock_requirement = self._mock_requirement(data.source_code)
        critical_conditions = self._extract_critical_conditions(
            source_code=data.source_code,
            language=data.language,
        )

        dependency_text = ", ".join(dependencies) if dependencies else "none detected"
        scenario_text = scenario if scenario else "general"
        complexity_text = complexity if complexity else "unknown"
        if critical_conditions:
            critical_text = "\n".join(f"- `{item}`" for item in critical_conditions)
        else:
            critical_text = "- none detected"

        return (
            "You are an expert unit testing engineer.\n"
            "你是一名资深单元测试工程师。\n"
            "Generate high-quality unit tests based on the following specification.\n"
            "请基于以下规范生成高质量单元测试。\n\n"
            "## Role & Objective（角色与目标）\n"
            "- Goal: produce executable tests that match source behavior exactly.\n"
            "- 目标：生成可执行且与源码行为严格一致的测试。\n\n"
            "## Language & Framework（语言与框架）\n"
            f"- Language（语言）: {data.language}\n"
            f"- Test Framework（测试框架）: {framework}\n\n"
            "## Step-by-Step Workflow（分步流程）\n"
            "1) Identify callable symbols and input/output contracts from source.\n"
            "2) Build a test matrix: normal, boundary, and exception paths.\n"
            "3) Derive expected values only from implementation semantics.\n"
            "4) Write deterministic, runnable tests with clear assertions.\n"
            "5) Self-check syntax/imports/assertions before final output.\n\n"
            "## Test Requirements（测试要求）\n"
            "- Cover normal paths, boundary conditions, and error/exception behavior.\n"
            "- 覆盖正常路径、边界条件和异常行为。\n"
            "- Keep tests deterministic and runnable.\n"
            "- 保持测试可重复、可执行（避免随机性）。\n"
            "- Use clear assertions with meaningful expected values.\n"
            "- 使用清晰断言和有意义的期望值。\n"
            "- Test function names must start with `test_`.\n"
            "- 测试函数命名必须以 `test_` 开头。\n"
            "- Use plain `assert` and `pytest.raises` for failure paths.\n"
            "- 断言使用 `assert`，异常路径使用 `pytest.raises`。\n"
            "- Use function-based tests (e.g., `def test_xxx()`) instead of class-based.\n"
            "- 使用函数式测试（如 `def test_xxx()`），不要用 class-based（如 `class Test:`）。\n"
            f"- MUST import target symbols from local module `{module_name}` before writing tests.\n"
            f"- 必须先从同目录模块 `{module_name}` 导入被测对象，再编写测试。\n"
            "- Avoid importing unrelated third-party packages by default.\n"
            "- 默认不要引入无关第三方包。\n"
            f"- Coverage targets（覆盖率目标，供参考）: {coverage_targets}\n"
            f"- Mock requirements（Mock 要求）: {mock_requirement}\n\n"
            "## Semantic Alignment Hard Rules（语义对齐硬约束）\n"
            "- Derive expected values strictly from the given source code behavior.\n"
            "- 期望值必须严格依据给定源码行为推导，不要按题型常识脑补。\n"
            "- Respect exact comparison semantics in code (`<`, `<=`, `>`, `>=`, `==`).\n"
            "- 必须严格遵守源码比较符号语义（尤其阈值边界等于时）。\n"
            "- Include explicit boundary-equality assertions when threshold/limit checks exist.\n"
            "- 当存在阈值/边界判断时，必须包含“等于边界”的断言样例。\n"
            "- If implementation looks counter-intuitive, still assert implementation behavior.\n"
            "- 若实现与常识不一致，也必须以源码实现为准。\n"
            "- If docstring/comment conflicts with implementation, trust implementation.\n"
            "- 若注释/文档示例与实现冲突，以实现为准。\n"
            "- Do NOT assume implicit coercion not present in code (e.g., str->number).\n"
            "- 不要假设源码未实现的隐式转换（例如字符串自动转数字）。\n"
            "- For sliding-window or two-pointer counting algorithms: when the loop condition is\n"
            "  `while left < right and sorted[right] - sorted[left] >= threshold`,\n"
            "  a threshold of 0 with duplicate values produces a count of 0 — the window only\n"
            "  advances when `>` threshold, NOT when `==` threshold.\n"
            "- 对于滑窗/双指针计数算法：当循环条件是 `>= threshold` 时，\n"
            "  threshold=0 且有重复值的情况下会计数为 0 —— 只有 `>` threshold 时窗口才右移。\n"
            "- Trace through the algorithm by hand for boundary values before writing assertions.\n"
            "- 写断言前，必须手工推导一遍算法在边界值上的执行过程。\n"
            "- For functions returning structured results (e.g., namedtuple, dataclass): always assert\n"
            "  each field individually. Never assert the whole object equality without field-level checks.\n"
            "- 对于返回结构体结果的函数，必须逐字段断言，切勿直接做整体相等判断而不验证字段值。\n"
            "- For `nearest_pair` or similar: the returned index fields are `left_index=min(original_indices)`\n"
            "  and `right_index=max(original_indices)` — trace the sorted enumeration carefully.\n"
            "- 对于类似 `nearest_pair` 的函数：返回的索引字段是 `left_index=min(原始索引)` 和\n"
            "  `right_index=max(原始索引)`，必须仔细追踪排序后的枚举过程。\n\n"
            "## Critical Conditions Extracted（关键逻辑条件）\n"
            f"{critical_text}\n\n"
            "## Error Prevention Checklist（错误预防清单，仅内部执行）\n"
            "- No placeholder tests like `assert True`.\n"
            "- No assertions for behavior that cannot be inferred from source.\n"
            "- Ensure every referenced symbol exists in source imports/definitions.\n"
            "- Ensure generated file is directly runnable by the target test framework.\n\n"
            "## Context Information（上下文信息）\n"
            f"- Sample ID（样本ID）: {data.sample_id}\n"
            f"- Scenario（场景）: {scenario_text}\n"
            f"- Complexity（复杂度）: {complexity_text}\n"
            f"- Dependencies detected（检测到依赖）: {dependency_text}\n\n"
            "## Output Format（输出格式）\n"
            "- Return raw test code only (no Markdown fences).\n"
            "- 仅输出原始测试代码，不要 Markdown 代码块。\n"
            "- Do not include explanations.\n"
            "- 不要输出解释文字。\n\n"
            "## Source Code Under Test（被测源码）\n"
            f"```{data.language}\n{data.source_code}\n```"
        )

    def _module_import_name(self, sample_id: str) -> str:
        normalized = re.sub(r"[^a-zA-Z0-9_]", "_", sample_id).strip("_")
        if not normalized:
            return "solution"
        if re.match(r"^[0-9]", normalized):
            return f"sample_{normalized}"
        return normalized

    def _infer_language(self, sample_path: Path, default_language: str | None) -> str:
        if default_language:
            return default_language.lower()

        ext = sample_path.suffix.lower()
        ext_map = {
            ".py": "python",
            ".java": "java",
            ".go": "go",
            ".cpp": "cpp",
            ".cc": "cpp",
            ".cxx": "cpp",
            ".js": "javascript",
            ".ts": "javascript",
        }
        if ext in ext_map:
            return ext_map[ext]
        return sample_path.parent.name.lower()

    def _coverage_targets_text(self) -> str:
        def pct(name: str, default: float) -> str:
            value = self._thresholds.get(name, default)
            try:
                return f"{float(value) * 100:.0f}%"
            except Exception:
                return f"{default * 100:.0f}%"

        return (
            f"line >= {pct('line_coverage', 0.7)}, "
            f"branch >= {pct('branch_coverage', 0.6)}, "
            f"function >= {pct('function_coverage', 0.8)}"
        )

    def _parse_sample_meta(
        self,
        sample_id: str,
        sample_path: Path,
    ) -> tuple[str | None, str | None]:
        # Expected shape: <complexity>_<language>_<scenario>_<index>
        parts = sample_id.split("_")
        if len(parts) < 4:
            parent = sample_path.parent.name.lower()
            scenario = parent if parent and parent != sample_path.parent.parent.name.lower() else None
            return scenario, None
        complexity = parts[0].lower()
        scenario = "_".join(parts[2:-1]).lower()
        return scenario, complexity

    def _extract_dependencies(self, source_code: str, language: str) -> list[str]:
        deps: list[str] = []
        if language == "python":
            for m in re.finditer(
                r"^\s*(?:from\s+([a-zA-Z0-9_\.]+)\s+import|import\s+([a-zA-Z0-9_\.]+))",
                source_code,
                flags=re.MULTILINE,
            ):
                dep = m.group(1) or m.group(2)
                if dep:
                    deps.append(dep)
        elif language == "java":
            deps.extend(
                re.findall(r"^\s*import\s+([^;]+);", source_code, flags=re.MULTILINE)
            )
        elif language == "go":
            deps.extend(re.findall(r'^\s*import\s+"([^"]+)"', source_code, flags=re.MULTILINE))
            block = re.search(r"import\s*\((.*?)\)", source_code, flags=re.DOTALL)
            if block:
                deps.extend(re.findall(r'"([^"]+)"', block.group(1)))
        elif language == "cpp":
            deps.extend(re.findall(r'^\s*#include\s*[<"]([^>"]+)[>"]', source_code, flags=re.MULTILINE))

        # Keep order while deduplicating.
        seen: set[str] = set()
        out: list[str] = []
        for dep in deps:
            key = dep.strip()
            if key and key not in seen:
                seen.add(key)
                out.append(key)
        return out[:12]

    def _mock_requirement(self, source_code: str) -> str:
        lowered = source_code.lower()
        external_markers = [
            "http",
            "request",
            "socket",
            "open(",
            "file",
            "database",
            "sql",
            "redis",
            "grpc",
            "client",
            "os.environ",
            "subprocess",
        ]
        if any(marker in lowered for marker in external_markers):
            return (
                "Use mocks/stubs/fakes for external dependencies "
                "(network, file system, database, subprocess, environment)."
            )
        return "Mock only when necessary; avoid over-mocking pure functions."

    def _extract_critical_conditions(self, source_code: str, language: str) -> list[str]:
        if language != "python":
            return []

        candidates: list[str] = []
        loop_lines: dict[int, str] = {}
        for lineno, raw_line in enumerate(source_code.splitlines(), 1):
            line = raw_line.strip()
            if not line or line.startswith("#"):
                continue
            if any(op in line for op in ("<=", ">=", "==", "!=", "<", ">")):
                if line.startswith("if ") or line.startswith("while ") or line.startswith("return "):
                    candidates.append(line)
                elif "while " in source_code[:source_code.find(line)] and "count" in line:
                    candidates.append(line)
            if any(line.startswith(p) for p in ("while ", "for ")) and any(op in line for op in ("<=", ">=", "==", "!=", "<", ">")):
                loop_lines[lineno] = line

        if loop_lines and not candidates:
            candidates.extend(loop_lines.values())

        dedup: list[str] = []
        seen: set[str] = set()
        for item in candidates:
            if item in seen:
                continue
            seen.add(item)
            dedup.append(item)
        return dedup[:12]
