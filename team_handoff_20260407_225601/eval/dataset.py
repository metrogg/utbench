import os
import json
import ast
import re
from dataclasses import dataclass, field
from typing import List, Dict, Optional


@dataclass
class CodeFile:
    language: str
    category: str
    filename: str
    filepath: str
    code: str
    test_ref: str = ""


@dataclass
class EvalResult:
    language: str
    category: str
    filename: str
    generated_test: str
    scores: Dict[str, float]
    total_score: float
    model_name: str


@dataclass
class ModelEvalReport:
    model_name: str
    results: List[EvalResult] = field(default_factory=list)

    @property
    def avg_total_score(self) -> float:
        if not self.results:
            return 0.0
        return sum(r.total_score for r in self.results) / len(self.results)

    @property
    def avg_by_language(self) -> Dict[str, float]:
        lang_scores = {}
        lang_counts = {}
        for r in self.results:
            lang_scores[r.language] = lang_scores.get(r.language, 0) + r.total_score
            lang_counts[r.language] = lang_counts.get(r.language, 0) + 1
        return {lang: lang_scores[lang] / lang_counts[lang] for lang in lang_scores}

    @property
    def avg_by_category(self) -> Dict[str, float]:
        cat_scores = {}
        cat_counts = {}
        for r in self.results:
            cat_scores[r.category] = cat_scores.get(r.category, 0) + r.total_score
            cat_counts[r.category] = cat_counts.get(r.category, 0) + 1
        return {cat: cat_scores[cat] / cat_counts[cat] for cat in cat_scores}

    @property
    def avg_by_dimension(self) -> Dict[str, float]:
        if not self.results:
            return {}
        dims = self.results[0].scores.keys()
        dim_scores = {d: 0.0 for d in dims}
        for r in self.results:
            for d in dims:
                dim_scores[d] += r.scores.get(d, 0)
        return {d: dim_scores[d] / len(self.results) for d in dims}


def load_code_files(base_dir: str = "data/basic_tests") -> List[CodeFile]:
    files = []
    for lang in os.listdir(base_dir):
        lang_dir = os.path.join(base_dir, lang)
        if not os.path.isdir(lang_dir):
            continue
        for fname in os.listdir(lang_dir):
            fpath = os.path.join(lang_dir, fname)
            if not os.path.isfile(fpath):
                continue
            with open(fpath, 'r', encoding='utf-8') as f:
                code = f.read()
            clean_name = fname
            for ext in ['.py', '.java', '.cpp', '.go']:
                if clean_name.endswith(ext):
                    clean_name = clean_name[:-len(ext)]
                    break
            parts = clean_name.split('_')
            category = parts[0] if parts else "unknown"
            files.append(CodeFile(
                language=lang,
                category=category,
                filename=fname,
                filepath=fpath,
                code=code
            ))
    return sorted(files, key=lambda x: (x.language, x.category))


def load_test_refs(json_path: str = "benchmark.json") -> Dict[str, str]:
    with open(json_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
    test_refs = {}
    for sample in data.get("samples", []):
        task_id = sample.get("id", "")
        test_ref = sample.get("test_ref", "")
        if test_ref:
            test_refs[task_id] = test_ref
    return test_refs
