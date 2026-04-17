import ast
import sys
from collections import Counter
from dataclasses import asdict, dataclass


DEFAULT_ALLOWED_NON_STDLIB = {
    "typing_extensions",
}


def stdlib_modules() -> set[str]:
    modules = set(getattr(sys, "stdlib_module_names", set()))
    modules.update(DEFAULT_ALLOWED_NON_STDLIB)
    return modules


def top_level_module(module_name: str | None) -> str | None:
    if not module_name:
        return None
    return module_name.split(".", 1)[0]


@dataclass(frozen=True)
class ImportIssue:
    issue_type: str
    module: str
    lineno: int


@dataclass(frozen=True)
class SelfContainedCheckResult:
    is_self_contained: bool
    imports: list[str]
    issues: list[ImportIssue]

    def issue_counts(self) -> dict[str, int]:
        return dict(Counter(issue.issue_type for issue in self.issues))

    def to_dict(self) -> dict:
        payload = asdict(self)
        payload["issue_counts"] = self.issue_counts()
        return payload


def analyze_self_contained_python(
    tree: ast.AST,
    *,
    package: str | None = None,
    allow_modules: set[str] | None = None,
) -> SelfContainedCheckResult:
    allowed_modules = stdlib_modules()
    if allow_modules:
        allowed_modules.update(allow_modules)

    imports: list[str] = []
    issues: list[ImportIssue] = []

    for node in ast.walk(tree):
        if isinstance(node, ast.Import):
            for alias in node.names:
                module_name = alias.name
                imports.append(module_name)
                module_root = top_level_module(module_name)
                if package and module_root == package:
                    issues.append(ImportIssue("internal_package_import", module_name, node.lineno))
                elif module_root and module_root not in allowed_modules:
                    issues.append(ImportIssue("third_party_import", module_name, node.lineno))
        elif isinstance(node, ast.ImportFrom):
            module_name = node.module or ""
            rendered_module = f"{'.' * node.level}{module_name}"
            imports.append(rendered_module or ".")
            if node.level:
                issues.append(ImportIssue("relative_import", rendered_module or ".", node.lineno))
                continue

            module_root = top_level_module(module_name)
            if package and module_root == package:
                issues.append(ImportIssue("internal_package_import", module_name, node.lineno))
            elif module_root and module_root not in allowed_modules:
                issues.append(ImportIssue("third_party_import", module_name, node.lineno))

    return SelfContainedCheckResult(
        is_self_contained=not issues,
        imports=sorted(set(imports)),
        issues=issues,
    )
