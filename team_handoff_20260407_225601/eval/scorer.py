import ast
import re
from typing import Dict


def evaluate_python_test(code: str, generated_test: str) -> Dict[str, float]:
    return _generic_evaluate(code, generated_test, "python")


def evaluate_java_test(code: str, generated_test: str) -> Dict[str, float]:
    return _generic_evaluate(code, generated_test, "java")


def evaluate_cpp_test(code: str, generated_test: str) -> Dict[str, float]:
    return _generic_evaluate(code, generated_test, "cpp")


def evaluate_go_test(code: str, generated_test: str) -> Dict[str, float]:
    return _generic_evaluate(code, generated_test, "go")


def evaluate_javascript_test(code: str, generated_test: str) -> Dict[str, float]:
    return _generic_evaluate(code, generated_test, "javascript")


def _generic_evaluate(code: str, generated_test: str, language: str) -> Dict[str, float]:
    scores = {
        "syntax_validity": 0.0,
        "test_coverage": 0.0,
        "assertion_quality": 0.0,
        "edge_case_handling": 0.0,
        "code_structure": 0.0
    }

    test_code = _extract_code_block(generated_test)

    scores["syntax_validity"] = _check_syntax(test_code, language)
    scores["test_coverage"] = _check_test_coverage(test_code, code, language)
    scores["assertion_quality"] = _check_assertions(test_code, language)
    scores["edge_case_handling"] = _check_edge_cases(test_code, language)
    scores["code_structure"] = _check_structure(test_code, language)

    weights = {
        "syntax_validity": 0.25,
        "test_coverage": 0.25,
        "assertion_quality": 0.20,
        "edge_case_handling": 0.15,
        "code_structure": 0.15
    }
    total = sum(scores[k] * weights[k] for k in scores)

    scores = {k: round(v, 2) for k, v in scores.items()}
    scores["total"] = round(total, 2)

    return scores


def _extract_code_block(text: str) -> str:
    patterns = [
        r"```(?:python|java|cpp|c\+\+|go|c|javascript|js)\s*\n(.*?)```",
        r"```\s*\n(.*?)```",
    ]
    for pattern in patterns:
        match = re.search(pattern, text, re.DOTALL)
        if match:
            return match.group(1).strip()
    return text.strip()


def _check_syntax(code: str, language: str) -> float:
    if not code or len(code) < 10:
        return 0.0

    if language == "python":
        try:
            ast.parse(code)
            return 1.0
        except SyntaxError:
            return 0.3
    elif language == "java":
        if "class" in code or "Test" in code:
            if code.count("{") == code.count("}"):
                return 0.8
            return 0.5
        return 0.2
    elif language == "cpp":
        if "TEST" in code or "test" in code.lower():
            if code.count("{") == code.count("}"):
                return 0.8
            return 0.5
        return 0.2
    elif language == "go":
        if "func Test" in code or "testing.T" in code:
            if code.count("{") == code.count("}"):
                return 0.8
            return 0.5
        return 0.2
    elif language == "javascript":
        if "describe" in code or "it(" in code or "test(" in code or "expect" in code:
            if code.count("{") == code.count("}"):
                return 0.8
            return 0.5
        return 0.2
    return 0.3


def _check_test_coverage(test_code: str, source_code: str, language: str) -> float:
    if not test_code or not source_code:
        return 0.0

    if language == "python":
        func_pattern = r"def\s+(\w+)\s*\("
    elif language == "java":
        func_pattern = r"(?:public|private|protected)?\s*\w+\s+(\w+)\s*\("
    elif language == "go":
        func_pattern = r"func\s+(\w+)\s*\("
    else:
        func_pattern = r"(?:\w+\s+)?(\w+)\s*\("

    source_funcs = set(re.findall(func_pattern, source_code))
    source_funcs = {f for f in source_funcs if len(f) > 2 and f not in ("if", "for", "while", "return", "import")}

    if not source_funcs:
        return 0.5

    covered = 0
    for func in source_funcs:
        if func in test_code:
            covered += 1

    return round(covered / len(source_funcs), 2)


def _check_assertions(code: str, language: str) -> float:
    if not code:
        return 0.0

    assertion_keywords = {
        "python": ["assert", "assertEquals", "assertTrue", "assertFalse", "assertEqual", "assertRaises", "pytest.raises"],
        "java": ["assert", "assertEquals", "assertTrue", "assertFalse", "assertThat", "assertNotNull", "assertNull"],
        "cpp": ["EXPECT_", "ASSERT_", "EXPECT_EQ", "ASSERT_EQ", "EXPECT_TRUE", "ASSERT_TRUE", "EXPECT_FALSE"],
        "go": ["t.Errorf", "t.Fatalf", "t.Error", "t.Fatal", "assert.", "require."],
        "javascript": ["expect", "assert", "toBe", "toEqual", "toBeTruthy", "toBeFalsy", "toBeNull", "toBeUndefined", "toHaveBeenCalled"]
    }

    keywords = assertion_keywords.get(language, [])
    found = sum(1 for kw in keywords if kw in code)

    if found == 0:
        return 0.1
    elif found <= 2:
        return 0.5
    elif found <= 4:
        return 0.8
    else:
        return 1.0


def _check_edge_cases(code: str, language: str) -> float:
    if not code:
        return 0.0

    edge_indicators = [
        "empty", "null", "nil", "None", "0", "-1", "boundary",
        "edge", "corner", "invalid", "error", "exception",
        "min", "max", "large", "small"
    ]

    code_lower = code.lower()
    found = sum(1 for ind in edge_indicators if ind in code_lower)

    if found == 0:
        return 0.1
    elif found <= 2:
        return 0.4
    elif found <= 4:
        return 0.7
    else:
        return 1.0


def _check_structure(code: str, language: str) -> float:
    if not code:
        return 0.0

    structure_score = 0.0

    if language == "python":
        if "class" in code and "Test" in code:
            structure_score += 0.4
        if "def test_" in code:
            structure_score += 0.3
        if "import" in code:
            structure_score += 0.3
    elif language == "java":
        if "class" in code and "Test" in code:
            structure_score += 0.4
        if "@Test" in code:
            structure_score += 0.3
        if "import" in code:
            structure_score += 0.3
    elif language == "cpp":
        if "TEST" in code:
            structure_score += 0.4
        if "EXPECT_" in code or "ASSERT_" in code:
            structure_score += 0.3
        if "#include" in code:
            structure_score += 0.3
    elif language == "go":
        if "func Test" in code:
            structure_score += 0.4
        if "testing.T" in code:
            structure_score += 0.3
        if "import" in code:
            structure_score += 0.3
    elif language == "javascript":
        if "describe" in code:
            structure_score += 0.4
        if "it(" in code or "test(" in code:
            structure_score += 0.3
        if "require" in code or "import" in code:
            structure_score += 0.3

    return round(structure_score, 2)
