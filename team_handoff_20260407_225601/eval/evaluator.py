import os
import json
import time
from datetime import datetime
from typing import List, Dict, Optional

from eval.dataset import CodeFile, EvalResult, ModelEvalReport, load_code_files
from eval.llm_client import LLMClient
from eval.scorer import (
    evaluate_python_test,
    evaluate_java_test,
    evaluate_cpp_test,
    evaluate_go_test,
    evaluate_javascript_test
)

SCORERS = {
    "python": evaluate_python_test,
    "java": evaluate_java_test,
    "cpp": evaluate_cpp_test,
    "go": evaluate_go_test,
    "javascript": evaluate_javascript_test
}


class UnitTestEvaluator:
    def __init__(self, api_key: str, model: str = "gpt-4o", output_dir: str = "eval/results", base_url: Optional[str] = None):
        self.client = LLMClient(api_key=api_key, model=model, base_url=base_url)
        self.model_name = model
        self.output_dir = output_dir
        os.makedirs(output_dir, exist_ok=True)

    def run_evaluation(self, code_files: List[CodeFile], save_results: bool = True) -> ModelEvalReport:
        report = ModelEvalReport(model_name=self.model_name)
        total = len(code_files)

        for i, cf in enumerate(code_files):
            print(f"[{i+1}/{total}] Evaluating {cf.language}/{cf.filename} ...")
            try:
                generated_test = self.client.generate_test(cf.code, cf.language, cf.category)
                scorer = SCORERS.get(cf.language)
                if scorer:
                    scores = scorer(cf.code, generated_test)
                else:
                    scores = {"syntax_validity": 0.5, "test_coverage": 0.5, "assertion_quality": 0.5,
                              "edge_case_handling": 0.5, "code_structure": 0.5, "total": 0.5}

                total_score = scores.get("total", 0.0)
                result = EvalResult(
                    language=cf.language,
                    category=cf.category,
                    filename=cf.filename,
                    generated_test=generated_test,
                    scores={k: v for k, v in scores.items() if k != "total"},
                    total_score=total_score,
                    model_name=self.model_name
                )
                report.results.append(result)
                print(f"  Score: {total_score:.2f}/1.00")

            except Exception as e:
                print(f"  ERROR: {e}")
                result = EvalResult(
                    language=cf.language,
                    category=cf.category,
                    filename=cf.filename,
                    generated_test="",
                    scores={"syntax_validity": 0, "test_coverage": 0, "assertion_quality": 0,
                            "edge_case_handling": 0, "code_structure": 0},
                    total_score=0.0,
                    model_name=self.model_name
                )
                report.results.append(result)

            time.sleep(1)

        if save_results:
            self._save_report(report)

        return report

    def _save_report(self, report: ModelEvalReport):
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"{self.model_name.replace('-', '_')}_{timestamp}.json"
        filepath = os.path.join(self.output_dir, filename)

        report_data = {
            "model_name": report.model_name,
            "total_files": len(report.results),
            "avg_total_score": round(report.avg_total_score, 4),
            "avg_by_language": {k: round(v, 4) for k, v in report.avg_by_language.items()},
            "avg_by_category": {k: round(v, 4) for k, v in report.avg_by_category.items()},
            "avg_by_dimension": {k: round(v, 4) for k, v in report.avg_by_dimension.items()},
            "results": [
                {
                    "language": r.language,
                    "category": r.category,
                    "filename": r.filename,
                    "total_score": r.total_score,
                    "scores": r.scores,
                    "generated_test": r.generated_test
                }
                for r in report.results
            ]
        }

        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(report_data, f, indent=2, ensure_ascii=False)

        print(f"\nReport saved to: {filepath}")

        md_path = os.path.join(self.output_dir, f"{self.model_name.replace('-', '_')}_{timestamp}.md")
        self._save_markdown_report(report, md_path)

    def _save_markdown_report(self, report: ModelEvalReport, filepath: str):
        lines = []
        lines.append(f"# Unit Test Generation Evaluation Report")
        lines.append(f"**Model**: {report.model_name}")
        lines.append(f"**Date**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        lines.append(f"**Total Files**: {len(report.results)}")
        lines.append("")

        lines.append("## Overall Score")
        lines.append(f"**Average Total Score**: {report.avg_total_score:.2f} / 1.00")
        lines.append("")

        lines.append("## Score by Language")
        lines.append("| Language | Avg Score |")
        lines.append("|----------|-----------|")
        for lang, score in sorted(report.avg_by_language.items()):
            lines.append(f"| {lang} | {score:.2f} |")
        lines.append("")

        lines.append("## Score by Category")
        lines.append("| Category | Avg Score |")
        lines.append("|----------|-----------|")
        for cat, score in sorted(report.avg_by_category.items()):
            lines.append(f"| {cat} | {score:.2f} |")
        lines.append("")

        lines.append("## Score by Dimension")
        lines.append("| Dimension | Avg Score |")
        lines.append("|-----------|-----------|")
        for dim, score in sorted(report.avg_by_dimension.items()):
            lines.append(f"| {dim} | {score:.2f} |")
        lines.append("")

        lines.append("## Detailed Results")
        lines.append("| Language | Category | File | Score |")
        lines.append("|----------|----------|------|-------|")
        for r in report.results:
            lines.append(f"| {r.language} | {r.category} | {r.filename} | {r.total_score:.2f} |")
        lines.append("")

        for r in report.results:
            lines.append(f"### {r.filename}")
            lines.append(f"- **Language**: {r.language}")
            lines.append(f"- **Category**: {r.category}")
            lines.append(f"- **Total Score**: {r.total_score:.2f}")
            lines.append(f"- **Dimension Scores**:")
            for dim, score in r.scores.items():
                lines.append(f"  - {dim}: {score:.2f}")
            lines.append("")
            lines.append("```")
            lines.append(r.generated_test[:2000])
            lines.append("```")
            lines.append("")

        with open(filepath, 'w', encoding='utf-8') as f:
            f.write('\n'.join(lines))

        print(f"Markdown report saved to: {filepath}")


def print_report_summary(report: ModelEvalReport):
    print("\n" + "=" * 60)
    print(f"EVALUATION REPORT: {report.model_name}")
    print("=" * 60)
    print(f"Total files evaluated: {len(report.results)}")
    print(f"Average total score: {report.avg_total_score:.2f} / 1.00")
    print()
    print("By Language:")
    for lang, score in sorted(report.avg_by_language.items()):
        print(f"  {lang:10s}: {score:.2f}")
    print()
    print("By Category:")
    for cat, score in sorted(report.avg_by_category.items()):
        print(f"  {cat:20s}: {score:.2f}")
    print()
    print("By Dimension:")
    for dim, score in sorted(report.avg_by_dimension.items()):
        print(f"  {dim:20s}: {score:.2f}")
    print()
    print("Per-file scores:")
    for r in report.results:
        print(f"  {r.language:6s} | {r.category:20s} | {r.filename:50s} | {r.total_score:.2f}")
    print("=" * 60)
